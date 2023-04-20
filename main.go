package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

// Deklarasi struktur mahasiswa
type Mahasiswa struct {
	ID                 int            `json:"id,omitempty"`
	Nama               string         `json:"nama,omitempty"`
	Usia               int            `json:"usia,omitempty"`
	Gender             int            `json:"gender,omitempty"`
	Tanggal_Registrasi sql.NullString `json:"tanggal_registrasi,omitempty"`
	ID_Jurusan         int            `json:"id_jurusan,omitempty"`
}

// Deklarasi struktur jurusan
type Jurusan struct {
	ID           int64  `json:"id,omitempty"`
	Nama_Jurusan string `json:"nama_jurusan,omitempty"`
}

// Deklarasi struktur hobi
type Hobi struct {
	ID        int64  `json:"id,omitempty"`
	Nama_Hobi string `json:"nama_hobi,omitempty"`
}

// Dekalarasi Struktur Mahasiswa Lengkap
type MahasiswaDetail struct {
	ID                 int      `json:"id,omitempty"`
	Nama               string   `json:"nama,omitempty"`
	Usia               int      `json:"usia,omitempty"`
	Gender             string   `json:"gender,omitempty"`
	Tanggal_Registrasi string   `json:"tanggal_registrasi,omitempty"`
	Nama_Jurusan       string   `json:"nama_jurusan,omitempty"`
	Nama_Hobi          []string `json:"nama_hobi,omitempty"`
}

// Deklarasi struktur mahasiswa versi gender String
type MahasiswaS struct {
	ID                 int    `json:"id,omitempty"`
	Nama               string `json:"nama,omitempty"`
	Usia               int    `json:"usia,omitempty"`
	Gender             string `json:"gender,omitempty"`
	Tanggal_Registrasi string `json:"tanggal_registrasi,omitempty"`
	ID_Jurusan         int    `json:"id_jurusan,omitempty"`
}

func main() {
	db, err = sql.Open("mysql", "root:@/techtestjh")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	router := mux.NewRouter()
	router.HandleFunc("/mahasiswas", getMahasiswas).Methods("GET")
	router.HandleFunc("/mahasiswa", createMahasiswa).Methods("POST")
	router.HandleFunc("/mahasiswa/{id}", getMahasiswa).Methods("GET")
	router.HandleFunc("/mahasiswa/{id}", updateMahasiswa).Methods("PUT")
	router.HandleFunc("/mahasiswa/{id}", deleteMahasiswa).Methods("DELETE")
	http.ListenAndServe(":8000", router)
}

func getMahasiswas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	results, err := db.Query("SELECT ID, Nama, Usia, Gender, Tanggal_Registrasi FROM Mahasiswa")
	if err != nil {
		panic(err.Error())
	}
	defer results.Close()

	var mahasiswas []MahasiswaS
	for results.Next() {
		var m Mahasiswa
		err := results.Scan(&m.ID, &m.Nama, &m.Usia, &m.Gender, &m.Tanggal_Registrasi)
		if err != nil {
			panic(err.Error())
		}
		var d MahasiswaS
		d.ID, d.Nama, d.Usia, d.Tanggal_Registrasi = m.ID, m.Nama, m.Usia, m.Tanggal_Registrasi.String
		if m.Gender == 0 {
			d.Gender = "Laki-Laki"
		} else {
			d.Gender = "Perempuan"
		}
		mahasiswas = append(mahasiswas, d)
	}
	if err := results.Err(); err != nil {
		panic(err.Error())
	}
	json.NewEncoder(w).Encode(mahasiswas)
}

func createMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	val := make(map[string]string)
	json.Unmarshal(body, &val)
	valNama := val["nama"]
	valUsia := val["usia"]
	valGender := val["gender"]
	valTanggal := val["tanggal_registrasi"]
	valJurusan := val["nama_jurusan"]
	valHobi := val["nama_hobi"]

	//Memeriksa nama jurusan apakah terdapt di DB atau tidak
	resultsJurusan, err := db.Query("Select ID, Nama_Jurusan FROM Jurusan")
	if err != nil {
		panic(err.Error())
	}
	defer resultsJurusan.Close()

	var jurusan []Jurusan
	for resultsJurusan.Next() {
		var nama Jurusan
		err := resultsJurusan.Scan(&nama.ID, &nama.Nama_Jurusan)
		if err != nil {
			panic(err.Error())
		}
		jurusan = append(jurusan, nama)
	}
	var idJurusan int64
	for i := 0; i < len(jurusan); i++ {
		if jurusan[i].Nama_Jurusan == valJurusan {
			idJurusan = jurusan[i].ID
			break
		} else if (i == len(jurusan)-1) && (jurusan[i].Nama_Jurusan != valJurusan) {
			stmtJ, err := db.Prepare("INSERT INTO Jurusan (Nama_Jurusan)VALUES (?)")
			if err != nil {
				panic(err.Error())
			}
			defer stmtJ.Close()
			result, err := stmtJ.Exec(valJurusan)
			if err != nil {
				panic(err.Error())
			}
			idJurusan, err = result.LastInsertId()
			if err != nil {
				panic(err.Error())
			}
		}
	}
	//Memeriksa nama hobi apakah terdapt di DB atau tidak
	resultsHobi, err := db.Query("Select ID, Nama_Hobi FROM Hobi")
	if err != nil {
		panic(err.Error())
	}
	defer resultsHobi.Close()

	var hobi []Hobi
	for resultsHobi.Next() {
		var nama Hobi
		err := resultsHobi.Scan(&nama.ID, &nama.Nama_Hobi)
		if err != nil {
			panic(err.Error())
		}
		hobi = append(hobi, nama)
	}
	var idHobi int64
	for i := 0; i < len(hobi); i++ {
		if hobi[i].Nama_Hobi == valHobi {
			idHobi = hobi[i].ID
			break
		} else if (i == len(hobi)-1) && (hobi[i].Nama_Hobi != valHobi) {
			stmtH, err := db.Prepare("INSERT INTO Hobi (Nama_Hobi)VALUES (?)")
			if err != nil {
				panic(err.Error())
			}
			defer stmtH.Close()
			result, err := stmtH.Exec(valHobi)
			if err != nil {
				panic(err.Error())
			}
			idHobi, err = result.LastInsertId()
			if err != nil {
				panic(err.Error())
			}
		}
	}

	//Memasukkan data mahasiswa
	stmt, err := db.Prepare("INSERT INTO Mahasiswa(Nama, Usia, Gender, Tanggal_Registrasi, ID_Jurusan) VALUES(?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	results, err := stmt.Exec(valNama, valUsia, valGender, valTanggal, idJurusan)
	mahasiswaID, err := results.LastInsertId()
	if err != nil {
		panic(err.Error())
	}

	//Memasukkan data Hobi di tabel Mahasiswa_Hobi
	insertMahasiswaHobi, err := db.Prepare("INSERT INTO Mahasiswa_Hobi(ID_Mahasiswa, ID_Hobi) VALUES(?,?)")
	if err != nil {
		panic(err.Error())
	}
	defer insertMahasiswaHobi.Close()
	res, err := insertMahasiswaHobi.Exec(mahasiswaID, idHobi)
	if err != nil {
		panic(err.Error())
	}

	json.NewEncoder(w).Encode(res)
}

func getMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	//Menarik data mahasiswa dari DB
	var mahasiswa Mahasiswa
	err := db.QueryRow("SELECT ID, Nama, Usia, Gender, Tanggal_Registrasi, ID_Jurusan FROM Mahasiswa WHERE ID = ?", params["id"]).Scan(&mahasiswa.ID, &mahasiswa.Nama, &mahasiswa.Usia, &mahasiswa.Gender, &mahasiswa.Tanggal_Registrasi, &mahasiswa.ID_Jurusan)
	if err != nil {
		panic(err.Error())
	}

	//Menarik data jurusan dari DB
	var jurusan Jurusan
	err = db.QueryRow("SELECT ID, Nama_Jurusan FROM Jurusan WHERE ID = ?", mahasiswa.ID_Jurusan).Scan(&jurusan.ID, &jurusan.Nama_Jurusan)
	if err != nil {
		panic(err.Error())
	}

	//Menarik data hobi dari DB
	results, err := db.Query("SELECT H.ID, H.Nama_Hobi FROM Hobi H INNER JOIN Mahasiswa_Hobi MH ON H.ID = MH.ID_Hobi WHERE MH.ID_Mahasiswa = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}

	defer results.Close()

	var hobi []Hobi
	for results.Next() {
		var h Hobi
		err := results.Scan(&h.ID, &h.Nama_Hobi)
		if err != nil {
			panic(err.Error())
		}
		hobi = append(hobi, h)
	}

	//Menjadikan semua data di atas menjadi format yang lebih mudah dibaca
	var data MahasiswaDetail
	data.ID, data.Nama, data.Usia, data.Tanggal_Registrasi = mahasiswa.ID, mahasiswa.Nama, mahasiswa.Usia, mahasiswa.Tanggal_Registrasi.String
	if mahasiswa.Gender == 0 {
		data.Gender = "Laki-Laki"
	} else {
		data.Gender = "Perempuan"
	}
	data.Nama_Jurusan = jurusan.Nama_Jurusan
	for i := 0; i < len(hobi); i++ {
		data.Nama_Hobi = append(data.Nama_Hobi, hobi[i].Nama_Hobi)
	}

	json.NewEncoder(w).Encode(data)
}

func updateMahasiswa(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	stmt, err := db.Prepare("UPDATE mahasiswa SET Nama = ?, Usia = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	val := make(map[string]string)
	json.Unmarshal(body, &val)
	valNama := val["nama"]
	valUsia := val["usia"]

	_, err = stmt.Exec(valNama, valUsia, params["id"])
	if err != nil {
		panic(err.Error())
	}
}

func deleteMahasiswa(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	//Melakukan delete data mahasiswa dari tabel Mahasiswa_Hobi
	stmt, err := db.Prepare("DELETE FROM Mahasiswa_Hobi WHERE ID_Mahasiswa = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}

	//Melakukan delete data mahasiswa dari tabel mahasiswa
	stmt, err = db.Prepare("DELETE FROM Mahasiswa WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}

}
