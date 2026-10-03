import 'package:finance_app/utils/meta.dart';
import 'package:flutter/material.dart';

class Themes {
  static Function(ThemeMode themeMode) changeThemeMode = (_) {};
  static String current = "system";

  static ThemeMode modeFromString(String mode) {
    switch (mode) {
      case "dark":
        return ThemeMode.dark;
      case "light":
        return ThemeMode.light;
      case "system":
      default:
        return ThemeMode.system;
    }
  }

  static String modeToString(ThemeMode mode) {
    switch (mode) {
      case ThemeMode.dark:
        return "dark";
      case ThemeMode.light:
        return "light";
      case ThemeMode.system:
        return "system";
    }
  }

  static ThemeData base(bool dark) {
    Brightness brightness = dark ? Brightness.dark : Brightness.light;
    return ThemeData(
      brightness: brightness,
      useMaterial3: false,
      appBarTheme: AppBarTheme(
          iconTheme: const IconThemeData(size: 19, color: Colors.white),
          toolbarHeight: 50,
          color: Meta.color[800],
          centerTitle: true),
      searchBarTheme: SearchBarThemeData(
          textStyle: MaterialStateProperty.resolveWith(
              (states) => const TextStyle(fontSize: 16))),
      searchViewTheme:
          const SearchViewThemeData(constraints: BoxConstraints(maxHeight: 50)),
      tabBarTheme: TabBarTheme(
        indicatorColor: Colors.grey[50],
      ),
      bottomNavigationBarTheme: BottomNavigationBarThemeData(
          backgroundColor: dark ? Colors.grey[850] : Colors.grey[50],
          selectedItemColor: dark ? Meta.color[500] : Meta.color[800]),
      scaffoldBackgroundColor: dark ? Colors.grey[850] : Colors.grey[50],
      colorScheme: ColorScheme.fromSwatch(
          primarySwatch: Meta.color,
          brightness: brightness,
          accentColor: Meta.color),
      fontFamily: "OpenSans",
    );
  }

  static final light = base(false);

  static final dark = base(true);

  static bool isDark(BuildContext context) {
    return Theme.of(context).brightness == Brightness.dark;
  }
}
