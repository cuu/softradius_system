<script>
  function deleteNode(node_id) {
  if (confirm("确认删除吗？")) {
  window.location.href = "/node/delete?node_id=" + node_id;
  }
  }
</script>

<div class="box">
  <div class="panel-heading"><i class="fa fa-map-marker"></i> 区域列表</div>
  <div class="panel-body">
    <div class="pull-right bottom10">
      {{ if eq (call .Match "/node/add") true }}
      <a href="/node/add" class="btn btn-sm btn-default">
	增加区域</a>
      {{end}}
    </div>
    <table class="table table-hover">
      <thead>
	<tr>
	  <th>区域名称</th>
	  <th>区域描述</th>
	  <th>操作</th>
	</tr>
      </thead>
      <tbody>
	{{range .Nodes }}
	<tr>
	  <td>{{.Name  }}</td>
	  <td>{{.Desc }}</td>
	  <td>
	    {{ if eq (call $.Match "/node/update") true }}
	    <a class="btn btn-default btn-xs" href="/node/update?node_id={{.Id}}">修改</a>
	    {{ end }}
	    
	    {{ if eq (call $.Match "/node/delete") true  }}
	    <a class="btn btn-default btn-xs btn-danger" href="javascript:deleteNode('{{.Id }}')">删除</a>
	    {{end}}
	  </td>
	</tr>
	{{end}}
      </tbody>
    </table>
    
  </div>
  </div> 
