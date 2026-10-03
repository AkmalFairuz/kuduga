import 'package:finance_app/constants.dart';
import 'package:finance_app/screens/splash_screen.dart';
import 'package:finance_app/themes.dart';
import 'package:finance_app/utils/local_storage.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/scroll/my_scroll_behavior.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

class App extends StatefulWidget {
  const App({super.key, this.initialAuthenticated = false});

  final bool initialAuthenticated;

  @override
  State<App> createState() => _AppState();
}

class _AppState extends State<App> {
  ThemeMode _themeMode = ThemeMode.system;

  @override
  void initState() {
    super.initState();

    SystemChrome.setPreferredOrientations([
      DeviceOrientation.portraitUp,
      DeviceOrientation.portraitDown,
    ]);

    Themes.changeThemeMode = _onThemeModeChange;

    _loadTheme();
  }

  void _onThemeModeChange(ThemeMode mode) async {
    setState(() {
      _themeMode = mode;
    });
    final themeModeStr = Themes.modeToString(mode);
    Themes.current = themeModeStr;
    await LocalStorage.set("themeMode", themeModeStr);
  }

  void _loadTheme() async {
    _themeMode = Themes.modeFromString(
        (await LocalStorage.get("themeMode", def: "system")) ?? "system");
    Themes.current = Themes.modeToString(_themeMode);
    setState(() {});
  }

  @override
  Widget build(BuildContext context) {
    return ScreenUtilInit(
      designSize: const Size(392, 805),
      minTextAdapt: true,
      builder: (_, child) {
        return MaterialApp(
          title: Constants.appName,
          localizationsDelegates: const [
            GlobalMaterialLocalizations.delegate,
            GlobalWidgetsLocalizations.delegate,
            GlobalCupertinoLocalizations.delegate,
          ],
          supportedLocales: const [Locale('id')],
          debugShowCheckedModeBanner: false,
          builder: (context, child) {
            return ScrollConfiguration(
                behavior: MyScrollBehavior(), child: child ?? Container());
          },
          navigatorKey: Screens.navigatorKey,
          darkTheme: Themes.dark,
          theme: Themes.light,
          themeMode: _themeMode,
          home: child,
        );
      },
      child: const SplashScreen(),
    );
  }
}
