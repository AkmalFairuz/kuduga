import 'package:flutter/material.dart';

class Spoiler extends StatefulWidget {
  const Spoiler({
    Key? key,
    required this.title,
    required this.child,
    this.initiallyExpanded = false,
    this.contentPadding =
        const EdgeInsets.only(left: 16, right: 16, bottom: 16, top: 2),
  }) : super(key: key);

  final Widget title;
  final Widget child;
  final bool initiallyExpanded;
  final EdgeInsets contentPadding;

  @override
  State<Spoiler> createState() => _SpoilerState();
}

class _SpoilerState extends State<Spoiler> {
  bool _isExpanded = false;

  @override
  void initState() {
    super.initState();
    _isExpanded = widget.initiallyExpanded;
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        InkWell(
          onTap: () {
            setState(() {
              _isExpanded = !_isExpanded;
            });
          },
          child: Container(
            padding: const EdgeInsets.symmetric(vertical: 8),
            child: Row(
              children: [
                const SizedBox(width: 16),
                Expanded(
                  child: widget.title,
                ),
                Icon(
                  _isExpanded
                      ? Icons.keyboard_arrow_up
                      : Icons.keyboard_arrow_down,
                  size: 24,
                ),
                const SizedBox(width: 8),
              ],
            ),
          ),
        ),
        if (_isExpanded)
          Padding(padding: widget.contentPadding, child: widget.child),
      ],
    );
  }
}
