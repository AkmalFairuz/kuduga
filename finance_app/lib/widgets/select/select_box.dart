import 'package:finance_app/controller/controller.dart';
import 'package:finance_app/utils/bottom_select.dart';
import 'package:finance_app/widgets/input/input.dart';
import 'package:flutter/material.dart';

class SelectController<T> extends Controller {
  T _value;

  SelectController(this._value);

  T value() {
    return _value;
  }

  void setValue(T value) {
    _value = value;
    notifyListeners();
  }
}

class SelectItem<T> {
  final T value;
  final String label;

  const SelectItem(this.value, this.label);
}

class SelectBox<T> extends StatefulWidget {
  const SelectBox(
      {Key? key,
      this.enabled = true,
      this.label,
      this.hintText,
      this.title,
      required this.items,
      this.onChange,
      this.controller})
      : super(key: key);

  final bool enabled;
  final String? hintText;
  final String? label;
  final String? title;
  final List<SelectItem<T>> items;
  final Function(T)? onChange;
  final SelectController<T>? controller;

  @override
  State<SelectBox> createState() => _SelectBoxState<T>();
}

class _SelectBoxState<T> extends State<SelectBox<T>> {
  var controller = SelectController<T?>(null);

  @override
  void initState() {
    super.initState();
    if (widget.controller != null) {
      controller = widget.controller!;
    } else {
      controller = SelectController<T?>(null);
    }
  }

  T? selected() {
    return controller.value();
  }

  @override
  Widget build(BuildContext context) {
    return Dyn(() {
      String currentLabel = selected() != null
          ? (widget.items
                  .where((element) => element.value == selected()!)
                  .toList()
                  .firstOrNull
                  ?.label ??
              "")
          : "";
      return Input(
        enabled: false,
        hintText: widget.hintText,
        label: widget.label,
        onTap: () {
          if (!widget.enabled) {
            return;
          }
          BottomSelect.show<T>(
                  context: context,
                  title: widget.label ?? widget.title ?? widget.hintText ?? "",
                  selectedValue: selected(),
                  items: widget.items
                      .map((e) => BottomSelectItem(
                          value: e.value, widget: Text(e.label)))
                      .toList())
              .then((v) {
            if (v != null) {
              setState(() {
                controller.setValue(v);
              });
              if (widget.controller != null) {
                widget.controller!.setValue(v);
              }
              if (widget.onChange != null) {
                (widget.onChange!)(v);
              }
            }
          });
        },
        controller: TextEditingController(text: currentLabel),
        suffixIcon: const Icon(Icons.arrow_drop_down_rounded),
      );
    }, widget.controller ?? SelectController(null));
  }
}
