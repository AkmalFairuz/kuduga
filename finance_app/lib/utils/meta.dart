import 'package:finance_app/screens/home_screen.dart';
import 'package:finance_app/screens/utility/use_referral_screen.dart';
import 'package:finance_app/utils/local_storage.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:intl/intl.dart';
import 'package:webview_flutter/webview_flutter.dart';

class Meta {
  static MaterialColor color = const MaterialColor(
    0xFF1976D2,
    {
      50: Color(0xFFE3F2FD),
      100: Color(0xFFBBDEFB),
      200: Color(0xFF90CAF9),
      300: Color(0xFF64B5F6),
      400: Color(0xFF42A5F5),
      500: Color(0xFF2196F3),
      600: Color(0xFF1E88E5),
      700: Color(0xFF1976D2),
      800: Color(0xFF1565C0),
      900: Color(0xFF0D47A1),
    },
  );

  static String? udid;

  static String currencyFormat(int value) {
    return value.toString().replaceAllMapped(
          RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'),
          (Match m) => '${m[1]}.',
        );
  }

  static String currencyFormatRp(int value) {
    if (value < 0) {
      return "-Rp${currencyFormat(value.abs())}";
    }
    return "Rp${currencyFormat(value)}";
  }

  static String formatUnixDate(int unixTimestamp,
      {String format = "d MMM yyyy HH:mm"}) {
    return DateFormat(format, 'id').format(
      DateTime.fromMillisecondsSinceEpoch(unixTimestamp * 1000),
    );
  }

  static String formatDateTime(DateTime dt,
      {String format = "d MMM yyyy HH:mm"}) {
    return DateFormat(format, 'id').format(dt);
  }

  static String buildQuery(Map<String, dynamic> params) {
    // convert all values to string
    params = params.map((key, value) => MapEntry(key, value.toString()));
    return Uri(queryParameters: params).query;
  }

  static Future<void> copyToClipboard(String text) {
    return Clipboard.setData(ClipboardData(text: text));
  }

  static void debug(dynamic data) {
    print(data);
  }

  static List<T> addSeparatorToList<T>(List<T> list, T separator,
      {bool addEnd = false}) {
    List<T> result = [];
    for (int i = 0; i < list.length; i++) {
      result.add(list[i]);
      if (i < list.length - 1) {
        result.add(separator);
      }
    }
    if (addEnd) {
      result.add(separator);
    }
    return result;
  }

  static bool isValidEmail(String email) {
    return RegExp(
            r"^[a-zA-Z0-9.a-zA-Z0-9.!#$%&'*+-/=?^_`{|}~]+@[a-zA-Z0-9]+\.[a-zA-Z]+")
        .hasMatch(email);
  }

  static bool isStringOnlyNumbers(String input) {
    final RegExp regex = RegExp(r'^[0-9]+$');
    return regex.hasMatch(input);
  }

  static bool isPin(String input) {
    return isStringOnlyNumbers(input) && input.length == 6;
  }

  static bool isImageExt(String file) {
    switch (file.split(".").last) {
      case "png":
      case "jpg":
      case "jpeg":
      case "bmp":
      case "webp":
      case "jfif":
        return true;
      default:
        return false;
    }
  }

  static String fileWithoutPath(String file) {
    return file.split("/").last;
  }

  static Future<void> webviewInjectFont(WebViewController controller) async {
    controller.runJavaScript("""
        var style = document.createElement('style');
        style.innerHTML = "@import url('https://fonts.googleapis.com/css2?family=Open+Sans:wght@400;700&display=swap');";
        document.head.appendChild(style);
        var style2 = document.createElement('style');
        style2.innerHTML = "body { font-family: 'Open Sans', sans-serif !important; } ";
        document.head.appendChild(style2);
      """);
  }

  static int unix() {
    return (DateTime.now().millisecondsSinceEpoch / 1000).round();
  }

  static void afterRegister() async {
    Screens.replaceAll(const HomeScreen());

    if (!((await LocalStorage.get("alreadyUsedReferral", def: false))!)) {
      Screens.to(const UseReferralScreen());
    }
  }
}
