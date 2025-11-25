import React from 'react';
import { Card, Button, Table, Alert, Row, Col } from 'react-bootstrap';
import { useCart } from './CartContext';

const CartPage = () => {
  const { cartItems, removeFromCart, cartTotal } = useCart();

  const handleCheckout = () => {
    alert(`Ukupna cena: €${cartTotal}. Nastavak na plaćanje...`);
    // TODO: Implementirati logiku za checkout
  };

  return (
    <div className="container-fluid mt-4">
      <Row>
        <Col>
          <Card>
            <Card.Header>
              <h3 className="mb-0">🛒 Moja korpa</h3>
            </Card.Header>
            <Card.Body>
              {cartItems.length === 0 ? (
                <Alert variant="info">
                  Vaša korpa je prazna. Idite na <a href="/tours">Dostupne ture</a> da dodate neku!
                </Alert>
              ) : (
                <>
                  <Table striped bordered hover responsive>
                    <thead>
                      <tr>
                        <th>#</th>
                        <th>Naziv Ture</th>
                        <th>Cena</th>
                        <th>Akcija</th>
                      </tr>
                    </thead>
                    <tbody>
                      {cartItems.map((item, index) => (
                        <tr key={item.tourId}>
                          <td>{index + 1}</td>
                          <td>{item.tourName}</td>
                          <td>€{parseFloat(item.price || 0).toFixed(2)}</td>
                          <td>
                            <Button 
                              variant="danger" 
                              size="sm" 
                              onClick={() => removeFromCart(item.id)}
                            >
                              Ukloni
                            </Button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </Table>

                  <div className="d-flex justify-content-end align-items-center mt-4">
                    <h4 className="me-3">Ukupno: <span className="text-success">€{cartTotal}</span></h4>
                    <Button 
                      variant="primary" 
                      size="lg" 
                      onClick={handleCheckout}
                    >
                      Nastavi sa plaćanjem
                    </Button>
                  </div>
                </>
              )}
            </Card.Body>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default CartPage;