package handler

import (
	"context"
	"net/http"
	"encoding/json"
	"errors"
	"strconv"
	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var req model.CreateTODORequest
		err:=json.NewDecoder(r.Body).Decode(&req)
		if err!=nil{
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if req.Subject==""{
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx:=r.Context()
		resp,err:=h.Create(ctx,&req)
		if err!=nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		err=json.NewEncoder(w).Encode(resp)
		if err!=nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}		
	}
	if r.Method == http.MethodPut {
		var req model.UpdateTODORequest
		err:=json.NewDecoder(r.Body).Decode(&req)
		if err!=nil{
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if req.ID==0 || req.Subject==""{
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx:=r.Context()
		resp,err:=h.Update(ctx,&req)
		if err!=nil{
			if errors.Is(err,&model.ErrNotFound{}){
				w.WriteHeader(http.StatusNotFound)
				return
			} else{
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		err=json.NewEncoder(w).Encode(resp)
		if err!=nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if r.Method == http.MethodGet {
		var req model.ReadTODORequest
		var err error

		val:=r.URL.Query()

		val_prev_id:=val.Get("prev_id")
		if val_prev_id==""{
			req.PrevID=0
		} else{
			req.PrevID,err=strconv.ParseInt(val_prev_id,10,64)
			if err!=nil{
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}


		val_size:=val.Get("size")
		if val_size==""{
			req.Size=0
		} else{
			req.Size,err=strconv.Atoi(val_size)
			if err!=nil{
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		ctx:=r.Context()
		resp,err:=h.Read(ctx,&req)
		if err!=nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err=json.NewEncoder(w).Encode(resp)
		if err!=nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if r.Method == http.MethodDelete {
		var req model.DeleteTODORequest
		err:=json.NewDecoder(r.Body).Decode(&req)
		if err!=nil{
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(req.IDs)==0{
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx:=r.Context()
		resp,err:=h.Delete(ctx,&req)

		if err!=nil{
			if errors.Is(err,&model.ErrNotFound{}){
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err=json.NewEncoder(w).Encode(resp)
		if err!=nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}

	var todo_res=&model.CreateTODOResponse{
		TODO: *todo,
	}
	return todo_res, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	todos, err := h.svc.ReadTODO(ctx, req.PrevID, int64(req.Size))
	if err != nil {
		return nil, err
	}

	var todo_res=&model.ReadTODOResponse{
		TODOs: todos,
	}

	return todo_res, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}

	var todo_res=&model.UpdateTODOResponse{
		TODO: *todo,
	}
	return todo_res, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	 err:= h.svc.DeleteTODO(ctx, req.IDs)
	 if err != nil {
		return nil, err
	 }
	return &model.DeleteTODOResponse{}, nil
}
