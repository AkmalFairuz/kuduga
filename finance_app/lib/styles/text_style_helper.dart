import 'package:finance_app/themes.dart';
import 'package:flutter/material.dart';

class TextStyleHelper {
  static TextStyle inputHelper(BuildContext context) {
    return TextStyle(
      color: Themes.isDark(context) ? Colors.grey[600] : Colors.grey[400],
      fontSize: 12,
    );
  }

  static TextStyle h1 = const TextStyle(
    fontSize: 24,
    fontWeight: FontWeight.bold,
  );
}
