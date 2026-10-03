import 'package:finance_app/utils/meta.dart';
import 'package:flutter/material.dart';

class ListButton extends StatelessWidget {
  const ListButton({
    Key? key,
    this.icon,
    required this.title,
    required this.onPress,
    this.suffix,
    this.suffixText,
  }) : super(key: key);

  final IconData? icon;
  final String title;
  final void Function() onPress;
  final Widget? suffix;
  final String? suffixText;

  @override
  Widget build(BuildContext context) {
    return InkWell(
        onTap: onPress,
        child: Container(
          padding: const EdgeInsets.symmetric(vertical: 20, horizontal: 16),
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.center,
            mainAxisAlignment: MainAxisAlignment.start,
            children: [
              if (icon != null) Icon(icon, size: 24, color: Meta.color[700]),
              const SizedBox(width: 16),
              Expanded(
                  child: Text(title, style: const TextStyle(fontSize: 15))),
              if (suffix != null) suffix!,
              if (suffix != null) const SizedBox(width: 8),
              if (suffixText != null)
                Text(suffixText!,
                    style: const TextStyle(fontSize: 12, color: Colors.grey)),
              if (suffixText != null) const SizedBox(width: 8),
              Icon(Icons.chevron_right, size: 24, color: Colors.grey[400]!),
            ],
          ),
        ));
  }
}
