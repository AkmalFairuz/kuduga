import 'package:finance_app/styles/color_helper.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

enum ButtonVariant { primary, outline, white }

class _ButtonVariant {
  final Color backgroundColor;
  final Color borderColor;
  final Color rippleColor;
  final double rippleOpacity;
  final TextStyle textStyle;

  const _ButtonVariant(
      {required this.backgroundColor,
      required this.borderColor,
      required this.rippleColor,
      required this.rippleOpacity,
      required this.textStyle});

  static _ButtonVariant get(BuildContext context, ButtonVariant variant) {
    switch (variant) {
      case ButtonVariant.primary:
        return _ButtonVariant(
            backgroundColor: ColorHelper.theme(context,
                dark: Meta.color[800], light: Meta.color[700])!,
            borderColor: Colors.transparent,
            rippleColor: ColorHelper.bg(context),
            rippleOpacity: 0.2,
            textStyle: const TextStyle(color: Colors.white));
      case ButtonVariant.outline:
        return _ButtonVariant(
            backgroundColor: ColorHelper.bg(context),
            borderColor: ColorHelper.theme(context,
                dark: Meta.color[800], light: Meta.color[700])!,
            rippleColor: Meta.color,
            rippleOpacity: 0.2,
            textStyle: TextStyle(
                color: ColorHelper.theme(context,
                    dark: Meta.color[500], light: Meta.color[700])!));
      case ButtonVariant.white:
        return _ButtonVariant(
            backgroundColor: Colors.white,
            borderColor: Colors.transparent,
            rippleColor: Colors.grey[800]!,
            rippleOpacity: 0.2,
            textStyle: TextStyle(color: Colors.grey[900]));
    }
  }
}

class Button extends StatelessWidget {
  // variant is optional, default to primary
  const Button({
    Key? key,
    this.variant = ButtonVariant.primary,
    required this.onPressed,
    required this.child,
    this.borderRadius = 4,
    this.width,
    this.height,
    this.padding = const EdgeInsets.symmetric(
      horizontal: 15,
      vertical: 9,
    ),
    this.isDisabled = false,
    this.isLoading = false,
  }) : super(key: key);

  final ButtonVariant variant;
  final VoidCallback onPressed;
  final Widget child;
  final double borderRadius;
  final double? width;
  final double? height;
  final EdgeInsets padding;
  final bool isLoading;
  final bool isDisabled;

  @override
  Widget build(BuildContext context) {
    _ButtonVariant variant = _ButtonVariant.get(context, this.variant);
    return Material(
      color: isDisabled
          ? variant.backgroundColor.withOpacity(0.25)
          : variant.backgroundColor,
      elevation: isDisabled ? 0 : 0.5,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(borderRadius),
        side: BorderSide(
          width: 1.5,
          color: isDisabled ? Colors.transparent : variant.borderColor,
        ),
      ),
      child: InkWell(
        onTap: isLoading || isDisabled ? null : onPressed,
        splashColor: variant.rippleColor.withOpacity(variant.rippleOpacity),
        splashFactory: InkRipple.splashFactory,
        borderRadius: BorderRadius.circular(borderRadius),
        child: Container(
          width: width,
          height: height,
          padding: padding,
          child: DefaultTextStyle(
            // use theme style then apply variant style
            textAlign: TextAlign.center,
            style: Theme.of(context)
                .textTheme
                .labelLarge!
                .merge(TextStyle(fontWeight: FontWeight.bold, fontSize: 15.sp))
                .merge(variant.textStyle),
            child: isLoading ? _loadingWidget() : child,
          ),
        ),
      ),
    );
  }

  Widget _loadingWidget() {
    return Row(mainAxisAlignment: MainAxisAlignment.center, children: [
      Container(
        width: 20.sp,
        height: 20.sp,
        padding: const EdgeInsets.symmetric(vertical: 1),
        child: CircularProgressIndicator(
          strokeWidth: 3.sp,
          valueColor: const AlwaysStoppedAnimation<Color>(Colors.white),
        ),
      ),
    ]);
  }
}
