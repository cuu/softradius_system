
<div class="panel-body">	
  {{if ne .msg  "" }}
  <div class="alert alert-warning">{{.msg}}</div>
  {{end}}
  
{{ str2html .Form }}
