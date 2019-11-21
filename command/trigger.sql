USE ginger_mall;

DROP TRIGGER IF EXISTS order_has_product_insert_trigger;
DELIMITER $$
CREATE TRIGGER order_has_product_insert_trigger AFTER INSERT ON ginger_mall.order_has_product FOR EACH ROW
BEGIN
    UPDATE product_pool_product SET sold_count=sold_count+new.count
    WHERE id=new.pool_product_id;
END
$$
DELIMITER ;
