import 'package:finance_app/widgets/layout/line.dart';
import 'package:flutter/material.dart';

import 'list_button.dart';

class ListButtonView extends StatelessWidget {
  const ListButtonView({
    Key? key,
    required this.buttons,
    this.separator = const Line(),
  }) : super(key: key);

  final List<ListButton> buttons;
  final Widget separator;

  @override
  Widget build(BuildContext context) {
    return ListView.separated(
      shrinkWrap: true,
      physics: const ClampingScrollPhysics(),
      itemCount: buttons.length,
      itemBuilder: (BuildContext context, int index) {
        return buttons[index];
      },
      separatorBuilder: (BuildContext context, int index) {
        return separator;
      },
    );
  }
}
