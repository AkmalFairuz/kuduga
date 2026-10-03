import 'package:finance_app/styles/color_helper.dart';
import 'package:flutter/material.dart';

class Fields {
  static List<Widget> build(BuildContext context, List<FieldEntry> entries,
      {bool striped = false}) {
    var i = 0;
    return entries.map((e) {
      i++;
      return buildSingle(context, e,
          color: striped
              ? (i % 2 != 0
                  ? ColorHelper.theme(context,
                      dark: Colors.grey[850], light: Colors.grey[100])
                  : ColorHelper.theme(context,
                      dark: Colors.grey[800], light: null))
              : null);
    }).toList();
  }

  static Widget buildSingle(BuildContext context, FieldEntry entry,
      {Color? color}) {
    return Container(
      color: color,
      padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 10),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(entry.key, style: TextStyle(color: ColorHelper.text(context))),
          Container(
              constraints: const BoxConstraints(
                maxWidth: 250,
              ),
              child: DefaultTextStyle(
                style: Theme.of(context).textTheme.bodyMedium!.copyWith(
                    color: ColorHelper.textLg(context),
                    fontWeight: FontWeight.w600),
                child: entry.value is String
                    ? SelectableText(
                        entry.value,
                        textAlign: TextAlign.end,
                      )
                    : entry.value,
              )),
        ],
      ),
    );
  }
}

class FieldEntry {
  String key;
  dynamic value;

  FieldEntry(this.key, this.value);
}
