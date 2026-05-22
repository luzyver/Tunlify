package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/luzyver/tunlify/internal/service"
)

type Projects struct {
	svc   *service.Projects
	audit *service.AuditLogger
}

func NewProjects(svc *service.Projects, audit *service.AuditLogger) *Projects {
	return &Projects{svc: svc, audit: audit}
}

func (h *Projects) List(w http.ResponseWriter, r *http.Request) {
	projects, err := h.svc.List()
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if projects == nil {
		projects = []service.Project{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func (h *Projects) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name        string `json:"name"`
		Path        string `json:"path"`
		RepoURL     string `json:"repo_url"`
		GitUsername string `json:"git_username"`
		GitToken    string `json:"git_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" || body.Path == "" {
		jsonError(w, "name and path required", http.StatusBadRequest)
		return
	}

	p, err := h.svc.Create(body.Name, body.Path, body.RepoURL, body.GitUsername, body.GitToken)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID, _ := r.Context().Value("user_id").(int)
	h.audit.Log(userID, "project_create", body.Name, r.RemoteAddr)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func (h *Projects) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var body struct {
		Name        string `json:"name"`
		Path        string `json:"path"`
		RepoURL     string `json:"repo_url"`
		GitUsername string `json:"git_username"`
		GitToken    string `json:"git_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" || body.Path == "" {
		jsonError(w, "name and path required", http.StatusBadRequest)
		return
	}

	if err := h.svc.Update(id, body.Name, body.Path, body.RepoURL, body.GitUsername, body.GitToken); err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID, _ := r.Context().Value("user_id").(int)
	h.audit.Log(userID, "project_update", body.Name, r.RemoteAddr)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Projects) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	if err := h.svc.Delete(id); err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID, _ := r.Context().Value("user_id").(int)
	h.audit.Log(userID, "project_delete", chi.URLParam(r, "id"), r.RemoteAddr)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Projects) History(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	p, err := h.svc.GetByID(id)
	if err != nil {
		jsonError(w, "project not found", http.StatusNotFound)
		return
	}
	logs, err := h.audit.ByDetail(p.Name, 20)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func (h *Projects) Output(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	lines, done := h.svc.GetOutput(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"lines": lines, "done": done})
}

func (h *Projects) Action(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	action := chi.URLParam(r, "action")

	p, err := h.svc.GetByID(id)
	if err != nil {
		jsonError(w, "project not found", http.StatusNotFound)
		return
	}

	var reqBody struct {
		Ref string `json:"ref"`
	}
	json.NewDecoder(r.Body).Decode(&reqBody)

	if action == "deploy" && reqBody.Ref == "" && p.RepoURL != "" {
		jsonError(w, "ref (branch/tag) required", http.StatusBadRequest)
		return
	}

	userID, _ := r.Context().Value("user_id").(int)
	ip := r.RemoteAddr

	onDone := func(output string) {
		detail := p.Name
		if action == "deploy" {
			detail = p.Name + "\n" + output
		}
		h.audit.Log(userID, "project_"+action, detail, ip)
	}

	switch action {
	case "up":
		h.svc.ComposeAction(id, p.Path, onDone, "up", "-d")
	case "down":
		h.svc.ComposeAction(id, p.Path, onDone, "down")
	case "restart":
		h.svc.ComposeAction(id, p.Path, onDone, "restart")
	case "deploy":
		h.svc.DeployStream(id, p, reqBody.Ref, onDone)
	default:
		jsonError(w, "invalid action", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
