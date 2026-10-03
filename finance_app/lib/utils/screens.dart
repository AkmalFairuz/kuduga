import 'package:flutter/material.dart';

class Screens {
  static GlobalKey<NavigatorState> navigatorKey = GlobalKey<NavigatorState>();

  static BuildContext _context() {
    return navigatorKey.currentContext!;
  }

  static Future<T?> to<T>(Widget widget) {
    return Navigator.push<T>(_context(), _pageRoute<T>(widget));
  }

  static void back() {
    Navigator.pop(_context());
  }

  static Future<T?> replace<T>(Widget widget) {
    return Navigator.pushReplacement(
      _context(),
      _pageRoute(widget),
    );
  }

  static void replaceAll(Widget widget) {
    Navigator.pushAndRemoveUntil(
      _context(),
      _pageRoute(widget),
      (route) => false,
    );
  }

  static PageRouteBuilder<T> _pageRoute<T>(Widget widget) {
    return PageRouteBuilder<T>(
      pageBuilder: (context, _, __) => widget,
      transitionDuration: const Duration(milliseconds: 200),
      transitionsBuilder: (context, animation, secondaryAnimation, child) {
        const begin = Offset(1.0, 0.0); // Start off-screen to the right
        const end = Offset.zero; // End at the center of the screen
        const curve = Curves.easeInOut; // Add your desired curve

        var tween =
            Tween(begin: begin, end: end).chain(CurveTween(curve: curve));

        var offsetAnimation = animation.drive(tween);

        return SlideTransition(
          position: offsetAnimation,
          child: child,
        );
      },
    );
  }

  static void unfocusInput() {
    var primaryFocus = FocusManager.instance.primaryFocus;
    if (primaryFocus != null) {
      primaryFocus.unfocus();
    }
  }
}
