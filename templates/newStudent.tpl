{{template "base" .}}
{{ define "embedded-css"}} {{end}}
{{define "title"}} Adicionar aluno {{end}}

{{define "header"}}
<div class="header">
    <h1>Adicionar aluno</h1>
    <h4></h4>
</div>
{{end}}

{{define "content"}}
<form class="content" action="/student/new" method="post">
    <label for="name">Nome completo</label>
    <input class="input" type="text" placeholder="" id="name" name="name" value={{.Name.Value}}>
    <p>{{.Name.Msg}}</p>

    <label for="email">E-mail</label>
    <input class="input" type="text" placeholder="" id="email" name="email" value={{.Email.Value}}>
    <p>{{.Email.Msg}}</p>

    <label for="mobile">Celular</label>
    <input class="input" type="text" placeholder="" id="mobile" name="mobile" value={{.Mobile.Value}}>
    <p>{{.Mobile.Msg}}</p>

    <button>Adicionar</button>
    <button>Cancelar</button>

</form>
{{end}}