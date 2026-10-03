import 'package:finance_app/screens/support/help_ticket_screen.dart';
import 'package:finance_app/service/support_service.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/layout/line.dart';
import 'package:finance_app/widgets/text/text_app_bar.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

import '../../model/support.dart';

class TicketListScreen extends StatefulWidget {
  const TicketListScreen({Key? key}) : super(key: key);

  @override
  State<StatefulWidget> createState() => _TicketListState();
}

class _TicketListState extends State<TicketListScreen> {
  List<SupportTicketModel> tickets = [];

  @override
  void initState() {
    super.initState();
    SupportService.getTickets().then((v) {
      setState(() {
        tickets = v;
      });
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const TextAppBar("Tiket Bantuan Anda"),
      body: ListView.separated(
          itemBuilder: (context, index) {
            var ticket = tickets[index];
            return InkWell(
                onTap: () => Screens.to(HelpTicketScreen(ticketId: ticket.id)),
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Row(
                    children: [
                      const Icon(CupertinoIcons.chat_bubble_2_fill),
                      const SizedBox(width: 16),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              ticket.category,
                              style: const TextStyle(fontSize: 16),
                              overflow: TextOverflow.ellipsis,
                            ),
                            const SizedBox(height: 2),
                            Row(
                              children: [
                                Text(
                                  "${ticket.statusText()} - ${Meta.formatUnixDate(ticket.createdAt)}",
                                  style: TextStyle(
                                      color: Colors.grey[500], fontSize: 12),
                                ),
                              ],
                            )
                          ],
                        ),
                      )
                    ],
                  ),
                ));
          },
          separatorBuilder: (_, __) => const Line(),
          itemCount: tickets.length),
    );
  }
}
