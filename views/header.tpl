
    <!-- Logo -->
    <a href="/" class="logo">
      <!-- mini logo for sidebar mini 50x50 pixels -->
      <span class="logo-mini"><b>S</b>R</span>
      <!-- logo for regular state and mobile devices -->
      <span class="logo-lg"><b>Soft</b>Radius</span>
    </a>

    <!-- Header Navbar: style can be found in header.less -->
    <nav class="navbar navbar-static-top">
      <!-- Sidebar toggle button-->
      <a href="#" class="sidebar-toggle" data-toggle="offcanvas" role="button">
        <span class="sr-only">Toggle navigation</span>
      </a>
      <!-- Navbar Right Menu -->
      <div class="navbar-custom-menu">
        <ul class="nav navbar-nav navbar-right">
	  <li >
	    <a href="/dashboard"><i class="fa fa-dashboard"></i>&nbsp;控制面板</a>
	  </li>
          <!-- User Account: style can be found in dropdown.less -->
          <li class="dropdown user user-menu">
            <a href="#" class="dropdown-toggle" data-toggle="dropdown">
              <img src="{{.adminlte}}/dist/img/android-icon-144x144.png" class="user-image" alt="User Image">
              <span class="hidden-xs">{{call .GetCookie "username"}}</span>
            </a>
            <ul class="dropdown-menu">
              <!-- User image -->
              <li class="user-header">
                <img src="{{.adminlte}}/dist/img/android-icon-144x144.png" class="img-circle" alt="User Image">
              </li>
              <!-- Menu Body -->
	      <!--
              <li class="user-body">
                <div class="row">
                  <div class="col-xs-4 text-center">
                    <a href="#">Followers</a>
                  </div>
                  <div class="col-xs-4 text-center">
                    <a href="#">Sales</a>
                  </div>
                  <div class="col-xs-4 text-center">
                    <a href="#">Friends</a>
                  </div>
                </div>
              </li>
	      -->
              <!-- Menu Footer-->
              <li class="user-footer">
		{{ if call $.Match "/opr/changepassword"}}
                <div class="pull-left">
                  <a href="/opr/changepassword" class="btn btn-default btn-flat">修改密码</a>
                </div>
		{{end}}
                <div class="pull-right">
                  <a href="/logout" class="btn btn-default btn-flat">Sign out</a>
                </div>
              </li>
            </ul>
          </li>
          <!-- Control Sidebar Toggle Button -->
        </ul>
      </div>

    </nav>
  </header>
