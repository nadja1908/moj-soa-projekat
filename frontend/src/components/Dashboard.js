import React from 'react';
import { Card, Row, Col, Badge } from 'react-bootstrap';
import { useAuth } from '../context/AuthContext';
import TailwindTest from './TailwindTest';

const Dashboard = () => {
  const { user } = useAuth();

  const getRoleBadgeVariant = (role) => {
    switch (role) {
      case 'administrator':
        return 'danger';
      case 'guide':
        return 'success';
      case 'tourist':
        return 'info';
      default:
        return 'secondary';
    }
  };

  const getRoleDisplayName = (role) => {
    switch (role) {
      case 'administrator':
        return 'Administrator';
      case 'guide':
        return 'Vodič';
      case 'tourist':
        return 'Turista';
      default:
        return role;
    }
  };

  const getRoleIcon = (role) => {
    switch (role) {
      case 'administrator':
        return '👑';
      case 'guide':
        return '🗺️';
      case 'tourist':
        return '🧳';
      default:
        return '👤';
    }
  };

  const getRoleDescription = (role) => {
    switch (role) {
      case 'administrator':
        return 'Imate punu kontrolu nad sistemom. Možete upravljati korisnicima i svim blog postovima.';
      case 'guide':
        return 'Možete kreirati blog postove o turističkim destinacijama i deliti svoja iskustva.';
      case 'tourist':
        return 'Možete čitati blog postove, ostavljati komentare i označavati omiljene postove.';
      default:
        return 'Dobrodošli na platformu!';
    }
  };

  const getAvailableActions = (role) => {
    const baseActions = [
      { icon: '📰', title: 'Čitanje blog postova', description: 'Pregledajte sve dostupne blog postove' },
      { icon: '💬', title: 'Komentarisanje', description: 'Ostavite komentare na blog postove' },
      { icon: '❤️', title: 'Označavanje omiljenih', description: 'Označite postove koji vam se sviđaju' }
    ];

    if (role === 'guide' || role === 'administrator') {
      baseActions.unshift({
        icon: '✍️', 
        title: 'Kreiranje blog postova', 
        description: 'Napišite i objavite svoje blog postove'
      });
    }

    if (role === 'administrator') {
      baseActions.push({
        icon: '⚙️', 
        title: 'Upravljanje korisnicima', 
        description: 'Blokiranje/odblokiranje korisnika'
      });
    }

    return baseActions;
  };

  return (
    <div>
      <div className="dashboard-header text-center">
        <div className="container">
          <h1 className="display-4 mb-3">🏠 Dashboard</h1>
          <p className="lead">Dobrodošli na SOA Blog Platform</p>
        </div>
      </div>

      <div className="container">
        <Row>
          <Col md={4}>
            <Card className="mb-4 border-0 shadow-sm">
              <Card.Body className="text-center">
                <div className="user-avatar mx-auto mb-3" style={{ width: '80px', height: '80px', fontSize: '2rem' }}>
                  {user.username.charAt(0).toUpperCase()}
                </div>
                <h4>{user.username}</h4>
                <Badge bg={getRoleBadgeVariant(user.role)} className="mb-2 role-badge">
                  {getRoleIcon(user.role)} {getRoleDisplayName(user.role)}
                </Badge>
                <p className="text-muted small">{user.email}</p>
                <p className="text-muted">{getRoleDescription(user.role)}</p>
              </Card.Body>
            </Card>
          </Col>

          <Col md={8}>
            <Card className="border-0 shadow-sm">
              <Card.Header className="bg-white border-0">
                <h5 className="mb-0">🎯 Dostupne funkcionalnosti</h5>
              </Card.Header>
              <Card.Body>
                <Row>
                  {getAvailableActions(user.role).map((action, index) => (
                    <Col md={6} key={index} className="mb-4">
                      <div className="d-flex align-items-start p-3 bg-light rounded">
                        <div className="feature-icon me-3">
                          {action.icon}
                        </div>
                        <div>
                          <h6 className="mb-2">{action.title}</h6>
                          <p className="text-muted small mb-0">{action.description}</p>
                        </div>
                      </div>
                    </Col>
                  ))}
                </Row>
              </Card.Body>
            </Card>

            <div className="admin-stats mt-4">
              <div className="admin-stat-card">
                <div className="admin-stat-number text-primary">-</div>
                <div className="admin-stat-label">Ukupno postova</div>
              </div>
              <div className="admin-stat-card">
                <div className="admin-stat-number text-success">-</div>
                <div className="admin-stat-label">Vaši komentari</div>
              </div>
              <div className="admin-stat-card">
                <div className="admin-stat-number text-danger">-</div>
                <div className="admin-stat-label">Označeni postovi</div>
              </div>
            </div>
            
            <div className="text-center mt-3">
              <small className="text-muted">
                Statistike će biti dostupne kada implementiramo dodatne funkcionalnosti
              </small>
            </div>
            
            {/* Tailwind CSS Test */}
            <div className="mt-5">
              <h5 className="mb-3">🎨 Tailwind CSS Test</h5>
              <TailwindTest />
            </div>
          </Col>
        </Row>
      </div>
    </div>
  );
};

export default Dashboard;