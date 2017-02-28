
<div class="user-panel">
  <div class="pull-left image">
    <img src="{{.adminlte}}/dist/img/android-icon-144x144.png" class="img-circle" alt="User Image">
  </div>
  <div class="pull-left info">
    <p>{{call  .GetCookie "username"}} </p>
    <a href="#"><i class="fa fa-circle text-success"></i> </a>
  </div>
</div>

<form action="/quicksearch" method="get" class="sidebar-form">
  <div class="input-group">
    <input type="text" name="q" class="form-control" placeholder="快捷搜索用户.">
    <span class="input-group-btn">
      <button type="submit" name="search" id="search-btn" class="btn btn-flat"><i class="fa fa-search"></i>
      </button>
    </span>
  </div>
</form>

<ul class="sidebar-menu">
  <li class="header">MAIN </li>
  <!-- -->

  {{ range .Menu }}
  {{ if (call $.CheckOprCate .Category)  }}
  <li class="treeview {{call $.Inactive . }}">
    <a href="#">
      <i class="{{index $.MenuIcon .Category }}"></i> <span>{{ .Category  }}</span>
      <span class="pull-right-container">
        <i class="fa fa-angle-left pull-right"></i>
      </span>
    </a>
    <ul class="treeview-menu">
      {{ range .Items }}
      {{ if and .Is_menu (call $.Match .Path) }}
      <li class="{{call $.AClass .Path}}" ><a href="{{ .Path }}"><i class="fa fa-circle-o"></i>{{ .Name }}</a></li>
      {{end}}
      {{ end }} 
    </ul>
  </li>
{{ end }}
  {{ end }}
  <!--
  <li><a href="{{.adminlte}}/documentation/index.html"><i class="fa fa-book"></i> <span>Documentation</span></a></li>
  
  <li class="header">LABELS</li>
  <li><a href="#"><i class="fa fa-circle-o text-red"></i> <span>Important</span></a></li>
  <li><a href="#"><i class="fa fa-circle-o text-yellow"></i> <span>Warning</span></a></li>
  <li><a href="#"><i class="fa fa-circle-o text-aqua"></i> <span>Information</span></a></li>
-->
</ul>
<div class="control-sidebar-bg"></div>
