import React, { useState, useEffect } from 'react';
import { Card, Table, Button, Badge, Alert, Spinner, Modal } from 'react-bootstrap';
import api from '../services/api';
import { useAuth } from '../context/AuthContext';

const AdminPanel = () => {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [actionLoading, setActionLoading] = useState(null);
  const [showModal, setShowModal] = useState(false);
  const [selectedUser, setSelectedUser] = useState(null);
  const { user } = useAuth();

  useEffect(() => {
    if (user && user.role === 'administrator') {
      fetchUsers();
    }
  }, [user]);

  const fetchUsers = async () => {
    try {
      setLoading(true);
      const response = await api.get('/users');
      setUsers(response.data || []);
    } catch (error) {
      console.error('Error fetching users:', error);
      setError('Greška pri učitavanju korisnika');
    } finally {
      setLoading(false);
    }
  };

  const handleBlockUser = async (userId, username) => {
    setSelectedUser({ id: userId, username });
    setShowModal(true);
  };

  const confirmBlockUser = async () => {
    if (!selectedUser) return;

    try {
      setActionLoading(selectedUser.id);
      await api.put(`/users/${selectedUser.id}/block`);
      await fetchUsers(); // Refresh the list
      setShowModal(false);
      setSelectedUser(null);
    } catch (error) {
      console.error('Error blocking user:', error);
      setError('Greška pri blokiranju korisnika');
    } finally {
      setActionLoading(null);
    }
  };

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

  if (!user || user.role !== 'administrator') {
    return (
      <Alert variant="danger" className="text-center">
        <Alert.Heading>⚠️ Nemate dozvolu</Alert.Heading>
        <p>Samo administratori mogu pristupiti ovoj stranici.</p>
      </Alert>
    );
  }

  if (loading) {
    return (
      <div className="text-center py-5">
        <Spinner animation="border" variant="primary" />
        <p className="mt-2 text-muted">Učitavanje korisnika...</p>
      </div>
    );
  }

  return (
    <div className="container">
      <div className="d-flex justify-content-between align-items-center mb-5">
        <div>
          <h1 className="display-5">⚙️ Admin Panel</h1>
          <p className="lead text-muted">Upravljanje korisnicima sistema</p>
        </div>
        <Badge bg="danger" className="p-3" style={{ fontSize: '1rem' }}>
          👑 Administrator
        </Badge>
      </div>

      {error && <Alert variant="danger" dismissible onClose={() => setError('')}>{error}</Alert>}

      <Card className="border-0 shadow-sm mb-4">
        <Card.Header className="bg-white border-0">
          <h5 className="mb-0">👥 Korisnici sistema ({users.length})</h5>
        </Card.Header>
        <Card.Body>
          {users.length === 0 ? (
            <div className="text-center py-5">
              <div className="feature-icon">👥</div>
              <p className="text-muted">Nema registrovanih korisnika.</p>
            </div>
          ) : (
            <Table responsive hover className="mb-0">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Korisnik</th>
                  <th>Email</th>
                  <th>Uloga</th>
                  <th>Status</th>
                  <th>Akcije</th>
                </tr>
              </thead>
              <tbody>
                {users.map((userItem) => (
                  <tr key={userItem.id}>
                    <td>{userItem.id}</td>
                    <td>
                      <div className="d-flex align-items-center">
                        <div className="user-avatar me-3" style={{ width: '35px', height: '35px', fontSize: '0.9rem' }}>
                          {userItem.username.charAt(0).toUpperCase()}
                        </div>
                        <strong>{userItem.username}</strong>
                      </div>
                    </td>
                    <td>{userItem.email}</td>
                    <td>
                      <Badge bg={getRoleBadgeVariant(userItem.role)} className="role-badge">
                        {getRoleIcon(userItem.role)} {getRoleDisplayName(userItem.role)}
                      </Badge>
                    </td>
                    <td>
                      <Badge bg={userItem.isActive ? 'success' : 'danger'}>
                        {userItem.isActive ? '✅ Aktivan' : '🚫 Blokiran'}
                      </Badge>
                    </td>
                    <td>
                      {userItem.role !== 'administrator' && userItem.id !== user.id && (
                        <Button
                          variant={userItem.isActive ? 'danger' : 'success'}
                          size="sm"
                          onClick={() => handleBlockUser(userItem.id, userItem.username)}
                          disabled={actionLoading === userItem.id}
                        >
                          {actionLoading === userItem.id ? (
                            <Spinner animation="border" size="sm" />
                          ) : userItem.isActive ? (
                            '🚫 Blokiraj'
                          ) : (
                            '✅ Odblokiraj'
                          )}
                        </Button>
                      )}
                      {(userItem.role === 'administrator' || userItem.id === user.id) && (
                        <small className="text-muted">-</small>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </Table>
          )}
        </Card.Body>
      </Card>

      <Card className="border-0 shadow-sm">
        <Card.Header className="bg-white border-0">
          <h5 className="mb-0">📊 Statistike sistema</h5>
        </Card.Header>
        <Card.Body>
          <div className="admin-stats">
            <div className="admin-stat-card">
              <div className="admin-stat-number text-primary">{users.length}</div>
              <div className="admin-stat-label">Ukupno korisnika</div>
            </div>
            <div className="admin-stat-card">
              <div className="admin-stat-number text-success">
                {users.filter(u => u.role === 'guide').length}
              </div>
              <div className="admin-stat-label">Vodiči</div>
            </div>
            <div className="admin-stat-card">
              <div className="admin-stat-number text-info">
                {users.filter(u => u.role === 'tourist').length}
              </div>
              <div className="admin-stat-label">Turisti</div>
            </div>
            <div className="admin-stat-card">
              <div className="admin-stat-number text-danger">
                {users.filter(u => u.role === 'administrator').length}
              </div>
              <div className="admin-stat-label">Administratori</div>
            </div>
          </div>
        </Card.Body>
      </Card>

      {/* Confirmation Modal */}
      <Modal show={showModal} onHide={() => setShowModal(false)}>
        <Modal.Header closeButton>
          <Modal.Title>⚠️ Potvrda akcije</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          {selectedUser && (
            <p>
              Da li ste sigurni da želite da {' '}
              <strong>
                {users.find(u => u.id === selectedUser.id)?.isActive ? 'blokirate' : 'odblokirate'}
              </strong>
              {' '} korisnika <strong>{selectedUser.username}</strong>?
            </p>
          )}
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={() => setShowModal(false)}>
            Odustani
          </Button>
          <Button 
            variant={users.find(u => u.id === selectedUser?.id)?.isActive ? 'danger' : 'success'} 
            onClick={confirmBlockUser}
            disabled={actionLoading}
          >
            {actionLoading ? (
              <Spinner animation="border" size="sm" />
            ) : users.find(u => u.id === selectedUser?.id)?.isActive ? (
              'Blokiraj'
            ) : (
              'Odblokiraj'
            )}
          </Button>
        </Modal.Footer>
      </Modal>
    </div>
  );
};

export default AdminPanel;