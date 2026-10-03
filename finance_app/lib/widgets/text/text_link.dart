import 'package:flutter/material.dart';

class TextLink extends StatelessWidget {
  const TextLink({
    Key? key,
    required this.text,
    required this.onPressed,
    this.textStyle,
  }) : super(key: key);

  final String text;
  final VoidCallback onPressed;
  final TextStyle? textStyle;

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onPressed,
      child: Text(
        text,
        // text style is bold and colored teal
        style: textStyle?.copyWith(
              fontWeight: FontWeight.bold,
              color: Colors.blue,
            ) ??
            const TextStyle(
              fontWeight: FontWeight.bold,
              color: Colors.blue,
            ),
      ),
    );
  }
}
