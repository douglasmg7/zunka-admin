{{ define "embedded-css"}}
<style type="text/css">
    .img-title {
        /* margin-bottom: .1em; */
        /* color: #3878BB; */
    }
    h2, h3 {
        margin-top: 0;
    }
    .panel {
        margin-top: 1em;
    }
    img {
        border-radius: 4px;
    }
</style>
{{end}}

{{template "base" .}}

{{define "title"}}Esola de vela Ventos Gerais{{end}}

{{define "content"}}
<div class="content">
    <div>
        <h1>Escola de Vela Ventos Gerais</h1>
        <h3>Seu primeiro grande passo no mundo da vela, seja ele para navegar outros oceanos ou como esporte.</h3>
    </div>
    <!-- Sub presentation. -->
    <div class="row">
        <div class="column panel">
            <img src="/static/img/sail_school.jpg" alt="Placeholder image">
            <h3 class="img-title">Aula para crianças</h4>
            <p>Criança também veleja, e os pequenos tem uma classe de vela só para eles: a Optimist!</p>
        </div>
        <div class="column panel">
            <img src="/static/img/sail3.jpg" alt="Placeholder image">
            <h3 class="img-title">Aula para adultos</h4>
            <p>Promovemos a iniciação de adultos na vela em diferentes barcos, para que cada aluno tenha uma rica experiência em diferentes embarcações.</p>
        </div>
        <div class="column panel">
            <img src="/static/img/sail_school2.jpg" alt="Placeholder image">
            <h3 class="img-title">Passeio e aluguel</h4>
            <p>É possível desfrutar da vela sem exatamente saber velejar. Nossa escola oferece passeios em veleiros na Lagoa dos Ingleses para até 6 pessoas, em barcos cabinados.</p>
        </div>
    </div>
</div>
{{end}}
