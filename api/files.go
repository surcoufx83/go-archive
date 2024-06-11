package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type File struct {
	ID                  uint           `json:"id"`
	DirectoryID         sql.NullInt32  `json:"directoryid"`
	ClassID             sql.NullInt16  `json:"classid"`
	CaseID              sql.NullInt16  `json:"caseid"`
	ClientParty         sql.NullInt16  `json:"clientparty"`
	PartyAddressID      sql.NullInt32  `json:"partyaddressid"`
	ContactID           sql.NullInt32  `json:"contactid"`
	NoClassify          uint8          `json:"noclassify"`
	Date                sql.NullString `json:"date"`
	Name                string         `json:"name"`
	FullPath            string         `json:"fullpath"`
	IsLink              uint8          `json:"islink"`
	LinkTo              sql.NullInt32  `json:"linkto"`
	MTime               string         `json:"mtime"`
	Size                uint64         `json:"size"`
	Hash                sql.NullString `json:"hash"`
	DelDate             sql.NullString `json:"deldate"`
	CaseFileName        sql.NullString `json:"case_filename"`
	CaseFileDescription sql.NullString `json:"case_filedescription"`
	CaseFileStatus      sql.NullString `json:"case_filestatus"`
	CasePinTop          sql.NullInt32  `json:"case_pintop"`
	FileClassMeta       sql.NullString `json:"fileclass_meta"`
	Updated             string         `json:"updated"`
	IsTaxReceipt        sql.NullInt32  `json:"istaxreceipt"`
	TaxYear             sql.NullInt16  `json:"taxyear"`
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Invalid file ID: %v", err)
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	var file File
	query := "SELECT id, directoryid, classid, caseid, clientparty, partyaddressid, contactid, noclassify, date, name, fullpath, islink, linkto, mtime, size, hash, deldate, case_filename, case_filedescription, case_filestatus, case_pintop, fileclass_meta, updated, istaxreceipt, taxyear FROM files WHERE id = ?"
	err = DB.QueryRow(query, id).Scan(
		&file.ID, &file.DirectoryID, &file.ClassID, &file.CaseID, &file.ClientParty, &file.PartyAddressID, &file.ContactID, &file.NoClassify, &file.Date, &file.Name, &file.FullPath, &file.IsLink, &file.LinkTo, &file.MTime, &file.Size, &file.Hash, &file.DelDate, &file.CaseFileName, &file.CaseFileDescription, &file.CaseFileStatus, &file.CasePinTop, &file.FileClassMeta, &file.Updated, &file.IsTaxReceipt, &file.TaxYear,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
		} else {
			log.Printf("Error querying database: %v", err)
			http.Error(w, "Error querying database", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(file); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
