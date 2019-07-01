{{template "base" .}}

{{define "title"}}Entrada{{end}}

{{define "body"}}
  <section class="section">
    <div class="container">

      <h1 class="title">Adicionar entrada</h2>

      <div class="field">
        <label class="label">Nome completo</label>
        <div class="control">
          <input class="input" type="text" placeholder="">
        </div>
      </div>

      <div class="field">
        <label class="label">Pulseira</label>
        <div class="control">
          <input class="input" type="text" placeholder="">
        </div>
      </div>

      <div class="field is-grouped">
        <div class="control">
          <button class="button is-link">Submit</button>
        </div>
        <div class="control">
          <button class="button is-text">Cancel</button>
        </div>
      </div>

    </div>
  </section>
{{end}}