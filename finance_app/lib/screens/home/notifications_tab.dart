import 'package:finance_app/service/notification_service.dart';
import 'package:finance_app/state/state_notifier.dart';
import 'package:finance_app/styles/color_helper.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/widgets/layout/line.dart';
import 'package:finance_app/widgets/select/horizontal_select.dart';
import 'package:finance_app/widgets/text/text_app_bar.dart';
import 'package:flutter/material.dart';

import '../../model/notification.dart';

class NotificationsTab extends StatefulWidget {
  const NotificationsTab({Key? key}) : super(key: key);

  @override
  State<NotificationsTab> createState() => _NotificationsTabState();
}

class _NotificationsTabState extends State<NotificationsTab> {
  List<NotificationModel>? _notifications;

  @override
  void initState() {
    super.initState();

    StateNotifier.getChannel("notification").listen(loadNotifications);
    loadNotifications();
  }

  @override
  void dispose() {
    StateNotifier.getChannel("notification").closeListener(loadNotifications);
    super.dispose();
  }

  Future<void> loadNotifications() async {
    var notifications = await NotificationService.getNotifications();
    setState(() {
      _notifications = notifications;
    });
  }

  Widget _buildContent() {
    if (_notifications != null) {
      return ListView.separated(
          itemBuilder: (context, index) {
            return _NotificationItem(_notifications![index]);
          },
          separatorBuilder: (_, __) => const Line(),
          itemCount: _notifications!.length);
    }
    return ListView(
      children: const [
        SizedBox(height: 40),
        Center(
          child: Text("Tidak ada notifikasi"),
        )
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        const TextAppBar("Notifikasi"),
        Padding(
            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
            child: HorizontalSelect(options: [
              HorizontalSelectOption(value: 0, widget: const Text("Semua"))
            ], selected: 0, onSelect: (_) {})),
        const Line(),
        Expanded(
            child: RefreshIndicator(
                onRefresh: loadNotifications, child: _buildContent()))
      ],
    );
  }
}

class _NotificationItem extends StatelessWidget {
  const _NotificationItem(this.model, {Key? key}) : super(key: key);

  final NotificationModel model;

  @override
  Widget build(BuildContext context) {
    return Material(
        color: ColorHelper.bg(context),
        child: InkWell(
            onTap: model.onClick,
            child: Container(
                padding:
                    const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Expanded(
                            child: Text(
                          model.title,
                          style: const TextStyle(fontWeight: FontWeight.bold),
                        )),
                        const SizedBox(width: 8),
                        Text(
                          Meta.formatUnixDate(model.createdAt),
                          style:
                              TextStyle(color: Colors.grey[500], fontSize: 12),
                        )
                      ],
                    ),
                    const SizedBox(height: 4),
                    Text(model.description)
                  ],
                ))));
  }
}
