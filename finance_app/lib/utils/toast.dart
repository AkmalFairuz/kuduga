import 'package:finance_app/utils/screens.dart';
import 'package:flutter/material.dart';

class ToastType {
  static const int success = 0;
  static const int error = 1;
  static const int warning = 2;
  static const int info = 3;
  static const int primary = 4;
}

class _ToastTypeInfo {
  final IconData icon;
  final Color color;
  final Color textColor;

  const _ToastTypeInfo({
    required this.icon,
    required this.color,
    required this.textColor,
  });
}

_ToastTypeInfo _getToastTypeInfo(int type) {
  switch (type) {
    case ToastType.success:
      return const _ToastTypeInfo(
        icon: Icons.check_circle_outline,
        color: Colors.green,
        textColor: Colors.white,
      );
    case ToastType.error:
      return const _ToastTypeInfo(
        icon: Icons.error_outline,
        color: Colors.red,
        textColor: Colors.white,
      );
    case ToastType.warning:
      return const _ToastTypeInfo(
        icon: Icons.warning_amber_outlined,
        color: Colors.orange,
        textColor: Colors.white,
      );
    case ToastType.primary:
      return _ToastTypeInfo(
        icon: Icons.info_outline,
        color: Colors.blue[700]!,
        textColor: Colors.white,
      );
  }
  return _ToastTypeInfo(
    icon: Icons.info_outline,
    color: Colors.white,
    textColor: Colors.grey[800]!,
  );
}

class Toast {
  static void primary(BuildContext context, String message, {String? subtext}) {
    text(message, type: ToastType.primary, subtext: subtext);
  }

  static void success(BuildContext context, String message, {String? subtext}) {
    text(message, type: ToastType.success, subtext: subtext);
  }

  static void error(BuildContext context, String message, {String? subtext}) {
    text(message, type: ToastType.error, subtext: subtext);
  }

  static void warning(BuildContext context, String message, {String? subtext}) {
    text(message, type: ToastType.warning, subtext: subtext);
  }

  static void info(String message, {String? subtext}) {
    text(message, type: ToastType.info, subtext: subtext);
  }

  static void text(String message,
      {int type = ToastType.info, String? subtext}) {
    _ToastTypeInfo info = _getToastTypeInfo(type);
    final ctx = Screens.navigatorKey.currentContext!;
    ScaffoldMessenger.of(ctx).showSnackBar(
      SnackBar(
          backgroundColor: info.color,
          behavior: SnackBarBehavior.floating,
          margin: const EdgeInsets.all(16),
          dismissDirection: DismissDirection.startToEnd,
          content: Row(
            children: [
              Icon(info.icon, color: info.textColor),
              const SizedBox(width: 20),
              Flexible(
                  child: Text(
                message,
                style: Theme.of(ctx).textTheme.bodyMedium!.copyWith(
                      color: info.textColor,
                      fontWeight: FontWeight.w600,
                    ),
              )),
            ],
          )),
    );
  }
}
