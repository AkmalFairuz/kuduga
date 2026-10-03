import 'package:flutter/material.dart';

class TextAppBar extends StatelessWidget implements PreferredSizeWidget {
  const TextAppBar(
    this.title, {
    Key? key,
  }) : super(key: key);

  final String title;

  @override
  Widget build(BuildContext context) {
    return AppBar(
      title: AppBarTitle(title),
    );
  }

  @override
  Size get preferredSize => const Size.fromHeight(50);
}

class AppBarTitle extends StatelessWidget {
  static const style = TextStyle(fontSize: 16, fontWeight: FontWeight.w700);

  const AppBarTitle(
    this.title, {
    Key? key,
  }) : super(key: key);

  final String title;

  @override
  Widget build(BuildContext context) {
    return Text(title, style: style);
  }
}
