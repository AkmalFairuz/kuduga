import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/utils/modal_bottom_sheet.dart';
import 'package:finance_app/widgets/text/modal_bottom_title.dart';
import 'package:flutter/material.dart';

class BottomSelect {
  static Future<T?> show<T>({
    required BuildContext context,
    required String title,
    required List<BottomSelectItem<T>> items,
    T? selectedValue,
  }) {
    //double screenHeight = MediaQuery.of(context).size.height;

    return ModalBottomSheet.show<T>(
        context,
        showDragHelper: false,
        (context) => Column(
              children: [
                ModalBottomTitle(title),
                const SizedBox(height: 8),
                _BottomSelectContent(
                    title: title, items: items, selectedValue: selectedValue)
              ],
            ));
  }
}

class _BottomSelectContent<T> extends StatefulWidget {
  const _BottomSelectContent({
    Key? key,
    required this.title,
    required this.items,
    required this.selectedValue,
  }) : super(key: key);

  final String title;
  final List<BottomSelectItem<T>> items;
  final T selectedValue;

  @override
  State<_BottomSelectContent> createState() => _BottomSelectContentState();
}

class _BottomSelectContentState<T> extends State<_BottomSelectContent> {
  T? _selectedValue;

  @override
  void initState() {
    super.initState();
    _selectedValue = widget.selectedValue;
  }

  @override
  Widget build(BuildContext context) {
    List<Widget> itemWidgets = widget.items.map((item) {
      return InkWell(
          onTap: () {
            setState(() {
              _selectedValue = item.value;
            });
            Navigator.pop(context, item.value);
          },
          child: Container(
            padding: const EdgeInsets.symmetric(vertical: 12),
            child: Row(
              children: [
                _selectedValue == item.value
                    ? Icon(Icons.check_circle, color: Meta.color[700], size: 20)
                    : const Icon(Icons.circle_outlined, size: 20),
                const SizedBox(width: 16),
                Expanded(
                    child: DefaultTextStyle(
                  // copy style from parent
                  style: Theme.of(context).textTheme.bodyMedium!.copyWith(
                        fontSize: 14.5,
                      ),
                  child: item.widget,
                )),
              ],
            ),
          ));
    }).toList();

    return Column(
      children: itemWidgets,
    );
  }
}

class BottomSelectItem<T> {
  final T value;
  final Widget widget;

  BottomSelectItem({required this.value, required this.widget});
}
