<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <title>{{.Title}}</title>
  <!-- Tell the browser to be responsive to screen width -->
  <meta content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no" name="viewport">
  <!-- Bootstrap 3.3.6 -->
  <link rel="stylesheet" href="{{.adminlte}}/bootstrap/css/bootstrap.min.css">
  <!-- Font Awesome -->
  <link rel="stylesheet" href="/static//css/font-awesome-4.7.0/css/font-awesome.min.css">
  <!-- <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.5.0/css/font-awesome.min.css"> -->
  <!-- Ionicons -->
  <link rel="stylesheet" href="/static/css/ionicons-2.0.1/css/ionicons.min.css">
<!--  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/ionicons/2.0.1/css/ionicons.min.css"> -->
  <!-- jvectormap -->
  <link rel="stylesheet" href="{{.adminlte}}/plugins/jvectormap/jquery-jvectormap-1.2.2.css">
  <!-- Theme style -->
  <link rel="stylesheet" href="{{.adminlte}}/dist/css/AdminLTE.css">
  <!-- AdminLTE Skins. Choose a skin from the css/skins
       folder instead of downloading all of them to reduce the load. -->
  <link rel="stylesheet" href="{{.adminlte}}/dist/css/skins/_all-skins.min.css">

  <link rel="stylesheet" href="/static/plugins/bootstrap-datepicker/dist/css/bootstrap-datepicker3.min.css" >
  
  <!-- HTML5 Shim and Respond.js IE8 support of HTML5 elements and media queries -->
  <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
  <!--[if lt IE 9]>
  <script src="https://oss.maxcdn.com/html5shiv/3.7.3/html5shiv.min.js"></script>
  <script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
  <![endif]-->

  <script src="{{.adminlte}}/plugins/jQuery/jquery-2.2.3.min.js"></script>
  <!-- Bootstrap 3.3.6 -->
  <script src="{{.adminlte}}/bootstrap/js/bootstrap.min.js"></script>
  <!-- FastClick -->
  <script src="{{.adminlte}}/plugins/fastclick/fastclick.js"></script>
  <!-- AdminLTE App -->
  <script src="{{.adminlte}}/dist/js/app.min.js"></script>
  <!-- Sparkline -->
  <script src="{{.adminlte}}/plugins/sparkline/jquery.sparkline.min.js"></script>
  <!-- jvectormap -->
  <script src="{{.adminlte}}/plugins/jvectormap/jquery-jvectormap-1.2.2.min.js"></script>
  <script src="{{.adminlte}}/plugins/jvectormap/jquery-jvectormap-world-mill-en.js"></script>
  <!-- SlimScroll 1.3.0 -->
  <script src="{{.adminlte}}/plugins/slimScroll/jquery.slimscroll.min.js"></script>
  <!-- ChartJS 1.0.1 -->
  <script src="{{.adminlte}}/plugins/chartjs/Chart.min.js"></script>

  <script src="/static/plugins/bootstrap-datepicker/dist/js/bootstrap-datepicker.min.js"></script>
  
  
  {{.HeadCss}}

  
</head>
<body class="hold-transition skin-blue sidebar-mini">
<div class="wrapper">
  <header class="main-header">
    {{.Header}}
  </header>

  <aside class="main-sidebar">
    <section class="sidebar">
      {{ .Sidebar}}
    </section>
  </aside>

  <div class="content-wrapper" >
    <section class="content-header">
      {{.ContentHeader}}
    </section>
    <section class="content">
      {{.LayoutContent}}
    </section>
  </div>

  <footer class="main-footer">
    {{.Footer}}
  </footer>
  
</div>
<!-- ./wrapper -->

  <!-- -------- -->
  
</body>
</html>
