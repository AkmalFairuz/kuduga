import 'package:finance_app/utils/screens.dart';
import 'package:flutter/material.dart';

class ModalBottomTitle extends StatelessWidget {
  const ModalBottomTitle(this.text,
      {Key? key, this.onClose, this.fontSize = 22})
      : super(key: key);

  final String text;
  final VoidCallback? onClose;
  final double fontSize;

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Expanded(
            child: Text(text,
                style: TextStyle(
                    fontWeight: FontWeight.bold, fontSize: fontSize))),
        const SizedBox(width: 16),
        GestureDetector(
          onTap: onClose ?? () => Screens.back(),
          child: const Icon(Icons.close, size: 24),
        )
      ],
    );
  }
}
