{{define "header"}}
    <a href="/">Home</a>
    <a href="/realtime">24H statistics</a>
    <a href="/history">Last statistics</a>
    <a href="/history-years">Global statistics</a>
    <nav class="navbar navbar-default navbar-static-top m-b-0">
            <div class="navbar-header">
                <div class="top-left-part">
                    <!-- Logo -->
                    <a class="logo" href="index.html">
                        <!-- Logo icon image, you can use font-icon also --><b>
                        <!--This is dark logo icon--><img src="../plugins/images/admin-logo.png" alt="home" class="dark-logo" /><!--This is light logo icon--><img src="../plugins/images/admin-logo-dark.png" alt="home" class="light-logo" />
                     </b>
                        <!-- Logo text image you can use text also --><span class="hidden-xs">
                        <!--This is dark logo text--><img src="../plugins/images/admin-text.png" alt="home" class="dark-logo" /><!--This is light logo text--><img src="../plugins/images/admin-text-dark.png" alt="home" class="light-logo" />
                     </span> </a>
                </div>
            </div>
            <!-- /.navbar-header -->
            <!-- /.navbar-top-links -->
            <!-- /.navbar-static-side -->
        </nav>
    <div class="navbar-default sidebar" role="navigation">
        <div class="sidebar-nav slimscrollsidebar">
            <div class="sidebar-head">
                <h3><span class="fa-fw open-close"><i class="ti-close ti-menu"></i></span> <span class="hide-menu">Menu</span></h3>
            </div>
            <ul class="nav" id="side-menu">
                <li style="padding: 70px 0 0;">
                    <a href="/" class="waves-effect"><i class="fa fa-clock-o fa-fw" aria-hidden="true"></i>Home</a>
                </li>
                <li>
                    <a href="/realtime" class="waves-effect"><i class="fa fa-user fa-fw" aria-hidden="true"></i>24H statistics</a>
                </li>
                <li>
                    <a href="/history" class="waves-effect"><i class="fa fa-table fa-fw" aria-hidden="true"></i>Last history</a>
                </li>
                <li>
                    <a href="/history-years" class="waves-effect"><i class="fa fa-font fa-fw" aria-hidden="true"></i>Global statistics</a>
                </li>
                <li>
                    <a href="/bots" class="waves-effect"><i class="fa fa-font fa-fw" aria-hidden="true"></i>Bots</a>
                </li>
            </ul>
        </div>
    </div>
{{end}}