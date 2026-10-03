import 'package:finance_app/utils/modal_bottom_sheet.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:flutter/material.dart';

class BottomButtons {
  static Future<void> show(BuildContext context, List<BottomButtonItem> items) {
    return ModalBottomSheet.show(context, (context) {
      return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: items.map((e) => e.build()).toList());
    });
  }
}

class BottomButtonItem {
  final String text;
  final VoidCallback onTap;

  const BottomButtonItem({required this.text, required this.onTap});

  Widget build() {
    return InkWell(
      onTap: () {
        Screens.back(); // close modal
        onTap();
      },
      child: Container(
          width: double.infinity,
          padding: const EdgeInsets.symmetric(vertical: 16, horizontal: 4),
          child: Text(
            text,
            style: const TextStyle(fontSize: 14.5),
          )),
    );
  }
}
