import 'package:finance_app/themes.dart';
import 'package:flutter/material.dart';

class ColorHelper {
  static Color bg(BuildContext context) {
    return Theme.of(context).scaffoldBackgroundColor;
  }

  static Color invertedBg(BuildContext context) {
    return invert(bg(context));
  }

  static Color bg0(BuildContext context) {
    return (Themes.isDark(context) ? Colors.grey[900] : Colors.grey[50])!;
  }

  static Color bg100(BuildContext context) {
    return (Themes.isDark(context) ? Colors.grey[800] : Colors.grey[100])!;
  }

  static Color bg200(BuildContext context) {
    return (Themes.isDark(context) ? Colors.grey[700] : Colors.grey[200])!;
  }

  static Color bg300(BuildContext context) {
    return (Themes.isDark(context) ? Colors.grey[600] : Colors.grey[300])!;
  }

  static Color bg400(BuildContext context) {
    return (Themes.isDark(context) ? Colors.grey[500] : Colors.grey[400])!;
  }

  static Color bg500(BuildContext context) {
    return (Themes.isDark(context) ? Colors.grey[400] : Colors.grey[500])!;
  }

  static Color bg600(BuildContext context) {
    return (Themes.isDark(context) ? Colors.grey[300] : Colors.grey[600])!;
  }

  static Color bg700(BuildContext context) {
    return (Themes.isDark(context) ? Colors.grey[200] : Colors.grey[700])!;
  }

  static Color bg800(BuildContext context) {
    return (Themes.isDark(context) ? Colors.grey[100] : Colors.grey[800])!;
  }

  static Color bg900(BuildContext context) {
    return (Themes.isDark(context) ? Colors.grey[50] : Colors.grey[900])!;
  }

  static Color mutedText(BuildContext context) {
    return (Themes.isDark(context) ? Colors.grey[700] : Colors.grey[400])!;
  }

  static Color textSm(BuildContext context) {
    return Theme.of(context).textTheme.bodySmall?.color ?? Colors.white;
  }

  static Color text(BuildContext context) {
    return Theme.of(context).textTheme.bodyMedium?.color ?? Colors.white;
  }

  static Color textLg(BuildContext context) {
    return Theme.of(context).textTheme.titleLarge?.color ?? Colors.white;
  }

  static Color invertedText(BuildContext context) {
    return invert(textSm(context));
  }

  static Color? theme(BuildContext context,
      {required Color? dark, required Color? light}) {
    return Themes.isDark(context) ? dark : light;
  }

  static Color invert(Color color) {
    return Color.fromARGB(
      color.alpha,
      255 - color.red,
      255 - color.green,
      255 - color.blue,
    );
  }

  static Color inputBg(BuildContext context) {
    return theme(context, dark: Colors.grey[900], light: Colors.grey[200])!;
  }
}
