import 'package:finance_app/styles/color_helper.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:flutter/material.dart';

class HorizontalSelect<T> extends StatelessWidget {
  const HorizontalSelect({
    Key? key,
    required this.options,
    required this.selected,
    required this.onSelect,
  }) : super(key: key);

  final List<HorizontalSelectOption> options;
  final T selected;
  final void Function(T) onSelect;

  @override
  Widget build(BuildContext context) {
    return SizedBox(
        height: 40,
        child: ListView.separated(
            scrollDirection: Axis.horizontal,
            itemCount: options.length,
            separatorBuilder: (context, index) => const SizedBox(width: 12),
            itemBuilder: (context, index) {
              return _HorizontalSelectOption(
                  selected: options[index].value == selected,
                  onPress: () {
                    onSelect(options[index].value);
                  },
                  child: options[index].widget);
            }));
  }
}

class HorizontalSelectOption<T> {
  HorizontalSelectOption({required this.value, required this.widget});

  final Widget widget;
  final T value;
}

class _HorizontalSelectOption extends StatelessWidget {
  const _HorizontalSelectOption({
    Key? key,
    required this.child,
    required this.selected,
    required this.onPress,
  }) : super(key: key);

  final Widget child;
  final bool selected;
  final void Function() onPress;

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onPress,
      child: Column(mainAxisAlignment: MainAxisAlignment.center, children: [
        AnimatedContainer(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
            decoration: BoxDecoration(
                color: selected
                    ? Meta.color[700]!
                    : ColorHelper.theme(context,
                        dark: Colors.grey[800], light: Colors.grey[200]),
                borderRadius: BorderRadius.circular(16)),
            duration: const Duration(milliseconds: 300),
            child: DefaultTextStyle(
                style: Theme.of(context).textTheme.bodyMedium!.copyWith(
                    color:
                        selected ? Colors.white : ColorHelper.textSm(context),
                    fontWeight: FontWeight.w600,
                    fontSize: 13),
                child: child))
      ]),
    );
  }
}
