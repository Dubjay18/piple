Grouped by resource, with role gating based on your `user_role` enum (`employee`, `procurement`, `ceo`, `admin`):

**Auth**
- `POST /auth/login`
- `POST /auth/refresh`
- `POST /auth/logout`
- No public `/auth/register` — accounts are admin-provisioned, not self-signup, since this is an internal payroll system.

**Users** (admin only)
- `GET /users` — filter by role, status
- `GET /users/:id`
- `POST /users`
- `PATCH /users/:id`
- `DELETE /users/:id` — should be soft delete (`deleted_at`), not hard, per the audit-trail decision earlier
- `GET /me` — current user's own profile

**Employees** (admin/HR write, self read-only on own record)
- `GET /employees` — filter by department, status, level
- `GET /employees/:id`
- `POST /employees`
- `PATCH /employees/:id` — includes bank detail updates, which should re-trigger Resolve Account Number before saving
- `DELETE /employees/:id`
- `GET /employees/:id/payouts` — that employee's payout history

**Salary codes** (admin only)
- `GET /salary-codes`
- `POST /salary-codes`
- `PATCH /salary-codes/:id`
- `DELETE /salary-codes/:id` — should be `restrict`/soft-delete if any `employees` reference it

**Payouts** (system-triggered + admin oversight)
- `GET /payouts` — filter by `user_id`, `status`, `pay_period`
- `GET /payouts/:id`
- `POST /payouts/run` — manual trigger for the batch (in addition to the 28th cron), useful for re-running a failed/partial month
- `POST /payouts/:id/retry` — re-attempt a single failed transfer without re-running the whole batch
- No public create — payouts are only ever system-generated from `salary_codes` + active employees, never POSTed ad-hoc

**Wallet** (CEO/admin only)
- `GET /wallet` — current balance
- `POST /wallet/topups/initialize` — kicks off Paystack Initialize Transaction, returns checkout URL/reference
- `GET /wallet/topups/:id` or `/wallet/topups/verify/:reference` — poll status as fallback to webhook
- `GET /wallet/topups` — history

**Payment requests** (procurement creates, CEO reviews)
- `POST /payment-requests` — procurement only
- `GET /payment-requests` — filter by status, `requested_by`; CEO sees all, procurement sees own
- `GET /payment-requests/:id`
- `PATCH /payment-requests/:id/approve` — CEO only, triggers transfer initiation
- `PATCH /payment-requests/:id/reject` — CEO only, should require a reason (add `rejection_reason` column if you don't have one — you don't currently)
- No `PATCH` for general edits once submitted — a payment request should be immutable except for status transitions, or you lose the audit trail you specifically asked to preserve

**Transactions** (read-only ledger)
- `GET /transactions` — filter by `wallet_id`, `user_id`, `type`, date range
- `GET /transactions/:id`
- No write endpoints — rows are only ever created by the payout/transfer/topup flows internally, never directly

**Banks** (thin proxy, no local table per your Paystack decision)
- `GET /banks` — proxies Paystack List Banks; cache in Redis with a short TTL since the list barely changes, don't hit Paystack on every page load
- `POST /banks/resolve-account` — proxies Resolve Account Number, `{account_number, bank_code}` → `{account_name}`

**Webhooks**
- `POST /webhooks/paystack` — single endpoint handling `charge.success`/`charge.failed` (wallet topups) and `transfer.success`/`transfer.failed`/`transfer.reversed` (payouts + payment requests). Verify the `x-paystack-signature` header against your secret key before processing — don't trust the payload otherwise.

