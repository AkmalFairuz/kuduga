import 'package:finance_app/styles/color_helper.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:flutter/material.dart';

class ModalBottomSheet {
  static Future<T?> show<T>(
    BuildContext context,
    WidgetBuilder builder, {
    EdgeInsets? padding,
    bool showDragHelper = true,
    bool isDismissible = true,
    bool enableDrag = true,
  }) {
    Screens.unfocusInput();
    return showModalBottomSheet<T>(
        isScrollControlled: true,
        // add top border radius
        shape: const RoundedRectangleBorder(
          borderRadius: BorderRadius.vertical(top: Radius.circular(12)),
        ),
        context: context,
        elevation: 2,
        isDismissible: isDismissible,
        enableDrag: enableDrag,
        builder: (context) {
          var bottomInner = MediaQuery.of(context).viewInsets.bottom;
          var screenHeight = MediaQuery.of(context).size.height;
          return Container(
              constraints: BoxConstraints(
                minHeight: bottomInner,
                maxHeight: screenHeight * 0.9,
              ),
              child: IntrinsicHeight(
                child: Column(
                  children: [
                    if (showDragHelper) ...[
                      const SizedBox(height: 12),
                      Container(
                        width: 30,
                        height: 5,
                        decoration: BoxDecoration(
                            color: ColorHelper.theme(context,
                                dark: Colors.grey[800],
                                light: Colors.grey[300]),
                            borderRadius: BorderRadius.circular(16)),
                      )
                    ],
                    Expanded(
                        child: SingleChildScrollView(
                      child: Container(
                        padding: (padding ?? const EdgeInsets.all(24)).copyWith(
                            top: showDragHelper ? 4 : null,
                            bottom: bottomInner + 16),
                        child: builder(context),
                      ),
                    ))
                  ],
                ),
              ));
        });
  }
}
