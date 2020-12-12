package store

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"../model"
	"github.com/lib/pq"
)

// UserRepository ...
type UserRepository struct {
	store *Store
}

// Create ...
func (r *UserRepository) Create(u *model.User) error {
	// if err := u.Validate(); err != nil {
	// 	fmt.Printf("%v\n", err)
	// 	return err
	// }

	if err := u.BeforeCreate(); err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	fmt.Printf("\n")
	return r.store.db.QueryRow(
		"INSERT INTO users (email, encrypted_password, name, role, \"group\") VALUES ($1, $2, $3, $4, $5) RETURNING id",
		u.Email,
		u.EncryptedPassword,
		u.Name,
		"",
		"",
	).Scan(&u.ID)
}

// Find ...
func (r *UserRepository) Find(id int) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, name, encrypted_password FROM users WHERE id = $1",
		id,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {

		return nil, err
	}

	return u, nil
}

// FindByEmail ...
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, name, encrypted_password, role, \"group\" FROM users WHERE email = $1",
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.Name,
		&u.EncryptedPassword,
		&u.Role,
		&u.Group,
	); err != nil {

		return nil, err
	}

	return u, nil
}

func (r *UserRepository) CheckTryState(email string) error {
	var tryState *bool

	if err := r.store.db.QueryRow("SELECT try_state FROM users WHERE email = $1", email).Scan(
		&tryState,
	); err != nil {
		return errors.New("This email not registered")
	}

	if tryState == nil {
		return nil
	}

	if *tryState == true {
		_, err := r.store.db.Exec("UPDATE users SET try_state=false WHERE email = $1", email)
		if err != nil {
			return err
		}
	}

	if *tryState == false {
		return errors.New("You are banned")
	}

	return nil
}

func (r *UserRepository) UpdateTryStateByID(id string, status *bool) error {
	_, err := r.store.db.Exec("UPDATE users SET try_state=$1 WHERE id = $2", status, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) UpdateTryStateByEmail(email string, status *bool) error {
	_, err := r.store.db.Exec("UPDATE users SET try_state=$1 WHERE email = $2", status, email)
	if err != nil {
		return err
	}

	return nil
}

// FindByEmailChamp ...
func (r *UserRepository) FindByEmailChamp(email, champ string) error {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT email FROM "+champ+"_group WHERE email = $1",
		email,
	).Scan(
		&u.Email,
	); err != nil {

		return err
	}

	return nil
}

// AddToChamp ...
func (r *UserRepository) AddToChamp(u *model.Champs, champ string) error {
	var group string

	_, err := r.store.db.Exec("DELETE FROM "+champ+"_group WHERE email = $1", u.Email)
	if err != nil {
		return err
	}

	_, err = r.store.db.Exec(
		"INSERT INTO "+champ+"_group (email, module, standnum, moduls, issue) VALUES ($1, $2, $3, $4, $5)",
		u.Email,
		u.Module,
		"",
		pq.Array(u.Moduls),
		false,
	)
	if err != nil {
		return err
	}

	if err := r.store.db.QueryRow("SELECT \"group\" FROM users WHERE email = $1", u.Email).Scan(
		&group,
	); err != nil {
		return err
	}

	g := strings.Split(group, ",")
	for i := range g {
		if g[i] == champ {
			return nil
		}
	}

	group += champ // group += "," + champ

	_, err = r.store.db.Exec(
		"UPDATE users  SET \"group\"=$1 WHERE email = $2",
		group,
		u.Email,
	)
	return err
}

func (r *UserRepository) AddToModule(champ, email, modules string) error {
	_, err := r.store.db.Exec(
		"UPDATE " + champ + "_stands SET module=$1 WHERE email = $2",
		modules,
		email,
	)

	return err
}

// GetChamp ...
func (r *UserRepository) GetChamp(group, email string) (*model.Champs, error) {
	s := &model.Champs{}
	q := "SELECT id, email, moduls, module, standnum, issue FROM " + group + "_group WHERE email = $1"
	if err := r.store.db.QueryRow(q, email).Scan(
		&s.ID,
		&s.Email,
		pq.Array(&s.Moduls),
		&s.Module,
		&s.Standnum,
		&s.Issue,
	); err != nil {
		return nil, err
	}

	return s, nil
}

// GetStand ...
func (r *UserRepository) GetStand(champ, email, module string) ([]model.Stand, error) {
	s := []model.Stand{}
	q := "SELECT id, datacenter, digi, address, exsi_user, exsi_pass, digi_user, digi_pass, email, module, port FROM " + champ + "_stands WHERE email = $1 and module = $2"
	rows, err := r.store.db.Query(q, email, module)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		t := model.Stand{}
		err := rows.Scan(
			&t.ID,
			&t.Datacenter,
			&t.Digi,
			&t.Address,
			&t.Esxiuser,
			&t.Esxipass,
			&t.Digiuser,
			&t.Digipass,
			&t.Email,
			&t.Module,
			&t.PortT,
		)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		s = append(s, t)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return s, nil
}

// IssueStand ...
func (r *UserRepository) IssueStand(group, email string) error {
	//
	q := "UPDATE " + group + "_stands SET email = $1 WHERE id IN ( SELECT id FROM " + group + "_stands WHERE standnum IN ( SELECT standnum FROM " + group + "_stands WHERE email = '' LIMIT 1)) "
	_, err := r.store.db.Exec(q, email)
	if err != nil {
		return err
	}
	//
	q = "UPDATE " + group + "_group SET issue = $1 WHERE email = $2"
	_, err = r.store.db.Exec(q, false, email)
	if err != nil {
		return err
	}
	return nil
}

// CreateTopology ...
func (r *UserRepository) CreateTopology(topology, name string) (int, error) {
	var id int

	err := r.store.db.QueryRow("INSERT INTO topologys (topology, name) VALUES ($1,$2) RETURNING id", topology, name).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// SaveTopology ...
func (r *UserRepository) SaveTopology(topology, name string) error {
	_, err := r.store.db.Exec("UPDATE topologys SET topology=$1 WHERE name = $2", topology, name)
	if err != nil {
		return err
	}

	return nil
}

// LinkTopology ...
func (r *UserRepository) LinkTopology(champ, module string, id int) error {
	var m string
	var t string
	var tid string

	err := r.store.db.QueryRow("SELECT moduls, topology FROM champs WHERE name = $1", champ).Scan(&m, &t)
	if err != nil {
		return err
	}
	modules := strings.Split(m, ",")
	topologys := strings.Split(t, ",")
	for i := range modules {
		if modules[i] == module {
			if i != 0 {
				tid = tid + "," + strconv.Itoa(id)
			} else {
				tid = tid + strconv.Itoa(id)
			}
		} else {
			if i != 0 {
				tid = tid + "," + topologys[i]
			} else {
				tid = tid + topologys[i]
			}
		}
	}
	fmt.Println(modules, topologys, module)
	_, err = r.store.db.Exec("UPDATE champs SET topology = $2, moduls = $3  WHERE name = $1", champ, tid, m)
	if err != nil {
		return err
	}

	err = r.store.db.QueryRow("SELECT champs FROM topologys WHERE id = $1", id).Scan(&m)
	if err != nil {
		return err
	}

	if m == "" {
		m = champ + ";" + module
	} else {
		m = m + "," + champ + ";" + module
	}

	_, err = r.store.db.Exec("UPDATE topologys SET champs = $2 WHERE id = $1", id, m)
	if err != nil {
		return err
	}

	return nil
}

// CreateChamp ...
func (r *UserRepository) CreateChamp(name string) error {
	var check int
	err := r.store.db.QueryRow("SELECT id FROM champs WHERE name = $1", name).Scan(&check)
	if err == nil {
		notfound := errors.New("Champ " + name + " is already exist")
		return notfound
	}

	_, err = r.store.db.Exec("INSERT INTO champs (name,datestart,dateend,moduls,topology) VALUES ($1,$2,$3,$4,$5);", name, nil, nil, "", "")
	if err != nil {
		return err
	}
	q := "CREATE Sequence " + name + "_group_id_seq;"
	_, err = r.store.db.Exec(q)
	if err != nil {
		return err
	}
	q = `CREATE TABLE public.` + name + `_group
	(
		id integer NOT NULL DEFAULT nextval('` + name + `_group_id_seq'::regclass),
		email character varying  NOT NULL,
		standnum character varying  NOT NULL,
		moduls character varying[]  NOT NULL,
		module character varying  NOT NULL,
		issue boolean NOT NULL DEFAULT true
	)`
	_, err = r.store.db.Exec(q)
	if err != nil {
		return err
	}
	_, err = r.store.db.Exec("CREATE Sequence " + name + "_stands_id_seq;")
	if err != nil {
		return err
	}
	q = `CREATE TABLE public.` + name + `;
	(
		id bigint NOT NULL DEFAULT nextval('` + name + `_stands_id_seq'::regclass),
		datacenter character varying  NOT NULL,
		digi character varying  NOT NULL,
		address character varying  NOT NULL,
		exsi_pass character varying  NOT NULL,
		exsi_user character varying  NOT NULL,
		digi_user character varying  NULL,
		digi_pass character varying  NOT NULL,
		email character varying  NOT NULL,
		module character varying  NOT NULL,
		port character varying  NOT NULL
	)
	`
	_, err = r.store.db.Exec(q)
	if err != nil {
		return err
	}
	return nil
}

// DeleteChamp ...
func (r *UserRepository) DeleteChamp(name string) error {
	_, err := r.store.db.Exec(
		"DELETE FROM champs where name = $1", name)
	if err != nil {
		return err
	}
	_, err = r.store.db.Exec("DROP TABLE public." + name + "_group;")
	if err != nil {
		return err
	}
	_, err = r.store.db.Exec("DROP Sequence " + name + "_group_id_seq;")
	if err != nil {
		return err
	}
	_, err = r.store.db.Exec("DROP TABLE public." + name + "_stands")
	if err != nil {
		return err
	}
	_, err = r.store.db.Exec("DROP Sequence " + name + "_stands_id_seq;")
	if err != nil {
		return err
	}
	return nil
}

// CreateModule ...
func (r *UserRepository) CreateModule(champ, module string) error {
	var modules, topology string
	fmt.Println(champ, module)
	err := r.store.db.QueryRow("SELECT moduls, topology FROM champs WHERE name = $1", champ).Scan(&modules, &topology)
	if err != nil {
		return err
	}
	if modules != "" {
		arr := strings.Split(modules, ",")
		for i := range arr {
			if arr[i] == module {
				return errors.New("Is already exist")
			}
		}
		modules = modules + "," + module
	} else {
		modules = module
	}
	topology = topology + ","

	_, err = r.store.db.Exec("UPDATE champs SET moduls = $2, topology = $3 WHERE name = $1", champ, modules, topology)
	if err != nil {
		return err
	}
	return nil
}

// DeleteModule ...
func (r *UserRepository) DeleteModule(champ, module string) error {
	var modules, topology string
	//var id int
	fmt.Println(champ, module)
	err := r.store.db.QueryRow("SELECT moduls, topology FROM champs WHERE name = $1", champ).Scan(&modules, &topology)
	if err != nil {
		return err
	}
	if modules != "" {
		arrM := strings.Split(modules, ",")
		arrT := strings.Split(topology, ",")
		for i := range arrM {
			if arrM[i] == module {
				copy(arrM[i:], arrM[i+1:])
				arrM[len(arrM)-1] = ""
				arrM = arrM[:len(arrM)-1]
				modules = strings.Join(arrM[:], ",")

				//id, err = strconv.Atoi(arrT[i])
				copy(arrT[i:], arrT[i+1:])
				arrT[len(arrT)-1] = ""
				arrT = arrT[:len(arrT)-1]
				topology = strings.Join(arrT[:], ",")
				break
			}
		}
	} else {
		return errors.New("Something went wrong")
	}
	fmt.Println(modules, topology)

	_, err = r.store.db.Exec("UPDATE champs SET moduls = $2, topology = $3 WHERE name = $1", champ, modules, topology)
	if err != nil {
		return err
	}
	// _, err = r.store.db.Exec("DELETE FROM topologys WHERE id = $1", id)
	// if err != nil {
	// 	return err
	// }
	return nil
}

// RemoveTopology ...
func (r *UserRepository) RemoveTopology(topology string) error {
	var champs, t, m string
	var id int
	err := r.store.db.QueryRow("DELETE FROM topologys WHERE name = $1 RETURNING champs, id", topology).Scan(&champs, &id)
	if err != nil {
		return err
	}
	if champs != "" {
		arr := strings.Split(champs, ",")
		fmt.Println(arr)
		for i := range arr {
			err := r.store.db.QueryRow("SELECT moduls,topology FROM champs WHERE name = $1", strings.Split(arr[i], ";")[0]).Scan(&m, &t)
			if err != nil {
				return err
			}

			if topology != "" {
				arrM := strings.Split(m, ",")
				arrT := strings.Split(t, ",")
				fmt.Println(arrM, arrT)
				for j := range arrM {
					tid := strconv.Itoa(id)
					fmt.Println(arrM, arrT, tid, strings.Split(arr[i], ";")[1])
					if arrM[j] == strings.Split(arr[j], ";")[1] && arrT[j] == tid {
						arrT[j] = ""
						t = strings.Join(arrT[:], ",")
						fmt.Println(t)
						break
					}
				}
			} else {
				continue
				//return errors.New("Something went wrong")
			}

			_, err = r.store.db.Exec("UPDATE champs SET topology = $2 WHERE name = $1", strings.Split(arr[i], ";")[0], t)
			if err != nil {
				return err
			}

		}
	}

	return nil
}

// GetTopology ...
func (r *UserRepository) GetTopology(name string) (int, string, error) {
	var tjson string
	var id int
	err := r.store.db.QueryRow("SELECT id, topology FROM topologys WHERE name = $1", name).Scan(&id, &tjson)
	if err != nil {
		return 0, "", err
	}
	return id, tjson, nil
}

// GetTopologyNameByID ...
func (r *UserRepository) GetTopologyNameByID(id int) (string, error) {
	var name string
	err := r.store.db.QueryRow("SELECT name FROM topologys WHERE id = $1", id).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

// GetTopologyForUser ...
func (r *UserRepository) GetTopologyForUser(champ, module string) (string, error) {
	var m, t string
	var tjson string
	var top string
	err := r.store.db.QueryRow("SELECT moduls, topology FROM champs WHERE name = $1", champ).Scan(&m, &t)
	if err != nil {
		return "", err
	}
	modules := strings.Split(m, ",")
	topologys := strings.Split(t, ",")
	for i := range modules {
		if modules[i] == module {
			top = topologys[i]
			break
		}
	}

	err = r.store.db.QueryRow("SELECT topology FROM topologys WHERE id = $1", top).Scan(&tjson)
	if err != nil {
		return "", err
	}
	return tjson, nil
}

// CreateStand ...
func (r *UserRepository) CreateStand(champ, module string) (int, error) {
	var id int
	q := "INSERT INTO " + champ + "_stands (digi, datacenter, address, exsi_user, exsi_pass, digi_user, digi_pass, email, port, module) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id "
	err := r.store.db.QueryRow(q, "", "", "", "", "", "", "", "", "", module).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// CreateStandCSV ...
func (r *UserRepository) CreateStandCSV(champ string, st *model.Stand) error {
	q := "INSERT INTO " + champ + "_stands (digi, address, exsi_user, exsi_pass, digi_user, digi_pass, module, port, datacenter, email) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)"
	_, err := r.store.db.Exec(q,
		&st.Digi,
		&st.Address,
		&st.Esxiuser,
		&st.Esxipass,
		&st.Digiuser,
		&st.Digipass,
		&st.Module,
		&st.PortT,
		&st.Datacenter,
		&st.Email,
	)
	if err != nil {
		return err
	}

	return nil
}

// RemoveStand ...
func (r *UserRepository) RemoveStand(champ string, id int) error {
	q := "DELETE FROM " + champ + "_stands WHERE id = $1"
	_, err := r.store.db.Exec(q, id)
	if err != nil {
		return err
	}
	return nil
}

// UpdateStand ...
func (r *UserRepository) UpdateStand(champ, module string, st *model.Stand) error {
	q := "UPDATE " + champ + "_stands SET digi=$1, address=$2, exsi_user=$3, exsi_pass=$4, digi_user=$5, digi_pass=$6, port=$8, datacenter=$9, email=$10 WHERE id = $7"
	_, err := r.store.db.Exec(q,
		&st.Digi,
		&st.Address,
		&st.Esxiuser,
		&st.Esxipass,
		&st.Digiuser,
		&st.Digipass,
		&st.ID,
		&st.PortT,
		&st.Datacenter,
		&st.Email,
	)
	if err != nil {
		return errors.New(q)
	}

	return nil
}

// AllStand ...
func (r *UserRepository) AllStand(champ, module string) ([](model.Stand), error) {
	stand := model.Stand{}
	stands := []model.Stand{}
	q := "SELECT digi, datacenter, address, exsi_user, exsi_pass, digi_user, digi_pass, port, id, email FROM " + champ + "_stands WHERE module = $1 ORDER BY id"
	fmt.Println("Query: ", q)
	rows, err := r.store.db.Query(q, module)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.Scan(
			&stand.Digi,
			&stand.Datacenter,
			&stand.Address,
			&stand.Esxiuser,
			&stand.Esxipass,
			&stand.Digiuser,
			&stand.Digipass,
			&stand.PortT,
			&stand.ID,
			&stand.Email,
		)
		if err != nil {
			log.Fatal(err)
			continue
			//return stands, err
		}
		stands = append(stands, stand)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return stands, nil
}

// UpdateTopology ...
func (r *UserRepository) UpdateTopology(top, champ, module string) error {
	var id int
	var m, t string
	err := r.store.db.QueryRow("SELECT moduls, topology FROM champs WHERE name = $1", champ).Scan(&m, &t)
	if err != nil {
		return err
	}
	modules := strings.Split(m, ",")
	topologys := strings.Split(t, ",")
	for i := range modules {
		if modules[i] == module {
			id, err = strconv.Atoi(topologys[i])
			break
		}
	}
	_, err = r.store.db.Exec("UPDATE topologys SET topology = $1 WHERE id = $2", top, id)
	if err != nil {
		return err
	}
	return nil
}

// GetModule ...
func (r *UserRepository) GetModule() ([]string, error) {
	var name string
	var s []string
	q := "SELECT name FROM topologys"
	rows, err := r.store.db.Query(q)
	if err != nil {
		return s, err
	}
	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
			return s, err
		}
		s = append(s, name)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
		return s, err
	}
	return s, nil
}

// AllChamp ...
func (r *UserRepository) AllChamp() ([]string, []string, error) {
	var a, b, c string
	var ch []string
	var ta []string

	q := "SELECT name, moduls, topology FROM champs"
	rows, err := r.store.db.Query(q)
	if err != nil {
		return nil, ch, err
	}
	for rows.Next() {
		err := rows.Scan(&a, &b, &c)
		if err != nil {
			log.Fatal(err)
			return nil, ch, err
		}
		ch = append(ch, a+","+b)
		ta = append(ta, c)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
		return nil, ch, err
	}
	return ta, ch, nil
}

// ResetPass ...
func (r *UserRepository) ResetPass(u *model.User) error {

	if err := u.BeforeCreate(); err != nil {
		fmt.Printf("%v\n", err)
		return err
	}

	q := "UPDATE users SET encrypted_password=$1 WHERE email = $2"
	_, err := r.store.db.Exec(q,
		u.EncryptedPassword,
		u.Email,
	)
	if err != nil {
		return err
	}

	return nil
}

// AllUser ...
func (r *UserRepository) AllUser() ([](model.User), error) {
	user := model.User{}
	users := []model.User{}
	q := "SELECT id, email, name, role, \"group\", try_state FROM users ORDER BY id"
	fmt.Println(q)
	rows, err := r.store.db.Query(q)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.Role,
			&user.Group,
			&user.TryState,
		)
		if err != nil {
			log.Fatal(err)
			continue
			//return stands, err
		}
		users = append(users, user)
	}
	// if err := rows.Err(); err != nil {
	// 	log.Fatal(err)
	// 	return nil, err
	// }

	return users, nil
}

// ChangeName ...
func (r *UserRepository) ChangeName(email, name string) error {
	q := "UPDATE users SET name=$1 WHERE email = $2"
	_, err := r.store.db.Exec(q,
		name,
		email,
	)
	if err != nil {
		return err
	}

	return nil
}
