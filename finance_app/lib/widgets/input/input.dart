import 'package:finance_app/styles/color_helper.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

class Input extends StatelessWidget {
  static TextStyle labelTextStyle(BuildContext context) {
    return TextStyle(
      color: ColorHelper.theme(context,
          dark: Colors.grey[200], light: Colors.grey[800]),
      fontWeight: FontWeight.bold,
      fontSize: 14,
    );
  }

  const Input({
    Key? key,
    this.focusNode,
    this.controller,
    this.hintText,
    this.isPassword = false,
    this.error,
    this.onTap,
    this.onChanged,
    this.onSubmitted,
    this.padding,
    this.borderRadius = 5,
    this.label,
    this.prefixIcon,
    this.suffixIcon,
    this.keyboardType,
    this.inputFormatters,
    this.enabled = true,
    this.minLines,
    this.maxLines = 1,
  }) : super(key: key);

  final VoidCallback? onTap;
  final FocusNode? focusNode;
  final bool enabled;
  final TextEditingController? controller;
  final String? hintText;
  final bool isPassword;
  final String? error;
  final ValueChanged<String>? onChanged;
  final ValueChanged<String>? onSubmitted;
  final EdgeInsetsGeometry? padding;
  final double borderRadius;
  final String? label;
  final Widget? prefixIcon;
  final Widget? suffixIcon;
  final TextInputType? keyboardType;
  final List<TextInputFormatter>? inputFormatters;
  final int? maxLines;
  final int? minLines;

  bool isError() {
    return error != null;
  }

  @override
  Widget build(BuildContext context) {
    return Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
      if (label != null) ...[
        Text(
          label!,
          style: labelTextStyle(context)
              .copyWith(color: isError() ? Colors.red[400] : null),
        ),
        const SizedBox(height: 4)
      ],
      Container(
        decoration: BoxDecoration(
          color: isError() ? Colors.red[50] : ColorHelper.inputBg(context),
          borderRadius: BorderRadius.circular(borderRadius),
        ),
        child: TextField(
          enableSuggestions: false,
          maxLines: maxLines,
          minLines: minLines,
          onTap: onTap,
          readOnly: !enabled,
          enableInteractiveSelection: enabled,
          focusNode: enabled ? focusNode : AlwaysDisabledFocusNode(),
          onTapOutside: (_) => Screens.unfocusInput(),
          keyboardType: keyboardType,
          controller: controller,
          obscureText: isPassword,
          onChanged: onChanged,
          onSubmitted: onSubmitted,
          style: TextStyle(
            color: isError()
                ? Colors.red[300]
                : ColorHelper.theme(context,
                    dark: Colors.grey[200], light: Colors.grey[800]),
            fontSize: 15,
          ),
          inputFormatters: inputFormatters,
          decoration: InputDecoration(
            suffixIcon: suffixIcon,
            suffixIconColor: isError() ? Colors.red[300] : Colors.grey[500],
            prefixIcon: prefixIcon,
            prefixIconColor: isError() ? Colors.red[300] : Colors.grey[500],
            contentPadding:
                padding ?? const EdgeInsets.symmetric(horizontal: 14),
            hintText: hintText,
            hintStyle: TextStyle(
              color: isError() ? Colors.red[300] : Colors.grey[500],
              fontSize: 15,
            ),
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(borderRadius),
              borderSide: BorderSide.none,
            ),
            focusedBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(borderRadius),
              borderSide: BorderSide(
                width: 2,
                color: isError() ? Colors.red : Meta.color[600]!,
              ),
            ),
          ),
        ),
      ),
      if (isError()) ...[
        const SizedBox(height: 2),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 4),
          child: Text(
            error!,
            style: TextStyle(
              color: Colors.red[400],
              fontSize: 12,
            ),
          ),
        ),
      ]
    ]);
  }
}

class AlwaysDisabledFocusNode extends FocusNode {
  @override
  bool get hasFocus => false;
}
