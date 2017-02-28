<script>
  function deleteOpr(opr_id)
  {
  if(confirm("确认删除吗？"))
  {
  window.location.href = "/opr/delete?opr_id="+opr_id;
  }
  }
  
</script>


<div class="box">
  <div class="panel-heading"><i class="fa fa-user-secret"></i> 操作员列表</div>
  <div class="panel-body">
    <div class="container">
      <div class="pull-right bottom10">
	{{ if eq (call .GetCookie "opr_type") "0" }}
	<a href="/opr/add" class="btn btn-sm btn-default">增加操作员</a>
	{{ end  }}
      </div>

      <table class="table table-hover">
	<thead>
	  <tr>
	    <th>操作员名称</th>
	    <th>操作员姓名</th>
	    <th>操作员类型</th>
	    <th>操作员状态</th>
	    <th>操作</th>
	  </tr>
	</thead>
	<tbody>
	  {{ range  .Opr_list }}
	  <tr>
	    <td>{{ .Name}}</td>
	    <td>{{ .Desc}}</td>
	    <td>{{map_get $.AllOprTypes .Type}}</td>
	    <td>{{map_get $.AllOprStatus .Status}}</td>
	    <td>
	      {{  if and (gt .Type  0) (eq (call $.GetCookie "opr_type") "0") }}
	      <a class="btn btn-default btn-xs"
		 href="/opr/update?opr_id={{.Id}}">修改</a>
	      <a class="btn btn-default btn-xs"
		 href="javascript:deleteOpr('{{.Id}}')">删除</a>
	      {{end}}
	    </td>
	  </tr>
	 {{end}}
	</tbody>
      </table>
    </div>
  </div>
</div>



