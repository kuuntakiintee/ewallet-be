### Prerequisites
### Database Setup

1.  Create a new database in PostgreSQL named `mywallet_db` (or whatever you prefer).
2.  Make sure your PostgreSQL server is running.

### Installation

1.  Install dependencies:
    ```bash
    go mod tidy
    ```

2.  **Environment Setup:**
    Create a `.env` file in the root directory and configure your credentials:

3.  Run the server:
    ```bash
    go run main.go
    ```
    *The application will automatically migrate the database schema upon startup.*

### API Endpoints

**Auth**
* `POST /api/register` - Create new account
* `POST /api/login` - Login user

**User**
* `GET /api/profile` - Get current user profile
* `GET /api/wallet` - Get wallet balance
* `POST /api/wallet/deposit` - Top up balance
* `POST /api/wallet/withdraw` - Withdraw balance
* `GET /api/transactions` - Get user history

**Admin**
* `GET /api/users` - List all users
* `POST /api/admin/users/:id/topup` - Inject balance
* `GET /api/admin/transactions` - Global history
