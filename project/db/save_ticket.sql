INSERT
INTO
	tickets
		(
			ticket_id,
			price_amount,
			price_currency,
			customer_email
		)
VALUES
	($1, $2, $3, $4);
