{{template "base" .}}

{{ define "embedded-css"}}
<style type="text/css">
    .content {
        max-width: var(--width-small);
    }
</style>
{{end}}

{{define "title"}}Recuperar senha{{end}}

{{define "content"}}
<form class="content" action="/auth/password/recovery" method="post">
    <h2 class="title">Recuperar senha</h2>

    <!-- email -->
    <label for="email">E-mail</label>
    <input type="text" id="email" name="email"  value={{.Email.Value}}>
    {{if .Email.Msg}} <p class="error"> {{.Email.Msg}} </p> {{end}}

    <!-- submit -->
    <input type="submit" value="Recuperar">

    <!-- Foot messages. -->
    {{if .SuccessMsgFooter}} <div class="success-msg"> {{.SuccessMsgFooter}} </div> {{end}}
    {{if .WarnMsgFooter}} <div class="warn-msg"> {{.WarnMsgFooter}} </div> {{end}}
</form>
{{end}}
