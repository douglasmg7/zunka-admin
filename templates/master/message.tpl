{{define "message"}}
{{ if .HeadMessage }}
<div class="alert alert-warning alert-dismissible fade show text-center mb-0" role="alert">
  <h5>{{.HeadMessage}}</h5>
  <button type="button" class="close" data-dismiss="alert" aria-label="Close">
    <span aria-hidden="true">&times;</span>
  </button>
</div>
{{end}}
{{end}}