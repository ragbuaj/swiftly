ALTER TABLE users 
    ADD COLUMN role VARCHAR(20) DEFAULT 'customer',
    ADD COLUMN status VARCHAR(20) DEFAULT 'active',
    ADD COLUMN email_verified_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN date_of_birth DATE,
    ADD COLUMN gender VARCHAR(10),
    ADD COLUMN newsletter_subscribed BOOLEAN DEFAULT false,
    ADD COLUMN avatar_url VARCHAR(255),
    ADD COLUMN bio TEXT,
    ADD COLUMN default_shipping_address_id UUID,
    ADD COLUMN default_billing_address_id UUID;
