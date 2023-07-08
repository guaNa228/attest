package main

// func (apiCfg *apiConfig) handlerCreateSemesterActivity(w http.ResponseWriter, r *http.Request, user db.User) {
// 	type parameters struct {
// 		GroupID   uuid.UUID `json:"group_id"`
// 		ClassID   uuid.UUID `json:"class_id"`
// 		TeacherID uuid.UUID `json:"teacher_id"`
// 	}
// 	decoder := json.NewDecoder(r.Body)
// 	params := parameters{}
// 	err := decoder.Decode(&params)
// 	if err != nil {
// 		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
// 		return
// 	}

// 	semesterActivityToCreate, err := apiCfg.DB.CreateSemesterActivity(r.Context(), db.CreateSemesterActivityParams{
// 		ID:        uuid.New(),
// 		CreatedAt: time.Now().UTC(),
// 		UpdatedAt: time.Now().UTC(),
// 		GroupID:   params.GroupID,
// 		ClassID:   params.ClassID,
// 		TeacherID: params.TeacherID,
// 	})

// 	if err != nil {
// 		respondWithError(w, 400, fmt.Sprintf("Couldn't create semester activity: %v", err))
// 		return
// 	}

// 	respondWithJSON(w, 201, semesterActivityToCreate)
// }

// func (apiCfg *apiConfig) handlerDeleteSemesterActivity(w http.ResponseWriter, r *http.Request, user db.User) {
// 	const instance = "semesterActivity"
// 	const paramToSearch = "semesterActivityToDeleteID"
// 	semesterActivityToDelete := chi.URLParam(r, paramToSearch)
// 	if semesterActivityToDelete == "" {
// 		respondWithError(w, 400, fmt.Sprintf("Wrong request address. Should be %v/{%vId}, not {%v}?{%vId}={%v}",
// 			instance,
// 			paramToSearch,
// 			instance,
// 			instance,
// 			paramToSearch))
// 	}

// 	semesterActivityToDeleteID, err := uuid.Parse(semesterActivityToDelete)

// 	if err != nil {
// 		respondWithError(w, 400, fmt.Sprintf("Corrupted %v id: %v", instance, err))
// 		return
// 	}

// 	err = apiCfg.DB.DeleteSemesterActivityById(r.Context(), semesterActivityToDeleteID)
// 	if err != nil {
// 		respondWithError(w, 400, fmt.Sprintf("Couldn't delete %v by ID: %v", instance, err))
// 		return
// 	}

// 	respondWithJSON(w, 200, struct{}{})
// }

// func (apiCfg *apiConfig) handlerUpdateSemesterActivity(w http.ResponseWriter, r *http.Request, user db.User) {
// 	const instance = "semesterActivity"
// 	const paramToSearch = "semesterActivityToUpdateID"
// 	semesterActivityToUpdate := chi.URLParam(r, paramToSearch)
// 	if semesterActivityToUpdate == "" {
// 		respondWithError(w, 400, fmt.Sprintf("Wrong request address. Should be {%v}/{%v}, not {%v}?{%v}Id={%v}",
// 			instance,
// 			paramToSearch,
// 			instance,
// 			instance,
// 			paramToSearch))
// 	}

// 	semesterActivityToUpdateID, err := uuid.Parse(semesterActivityToUpdate)

// 	if err != nil {
// 		respondWithError(w, 400, fmt.Sprintf("Corrupted %v id: %v", instance, err))
// 		return
// 	}

// 	type parameters struct {
// 		GroupID   uuid.UUID `json:"group_id"`
// 		ClassID   uuid.UUID `json:"class_id"`
// 		TeacherID uuid.UUID `json:"teacher_id"`
// 	}

// 	decoder := json.NewDecoder(r.Body)
// 	params := parameters{}
// 	err = decoder.Decode(&params)
// 	if err != nil {
// 		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
// 		return
// 	}

// 	newSemesterActivity, err := apiCfg.DB.UpdateSemesterActivityById(r.Context(), db.UpdateSemesterActivityByIdParams{
// 		ID:        semesterActivityToUpdateID,
// 		GroupID:   params.GroupID,
// 		TeacherID: params.TeacherID,
// 		ClassID:   params.ClassID,
// 	})
// 	if err != nil {
// 		respondWithError(w, 400, fmt.Sprintf("Couldn't update %v by ID: %v", instance, err))
// 		return
// 	}

// 	respondWithJSON(w, 200, newSemesterActivity)
// }
