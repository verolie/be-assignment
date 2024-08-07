-- CreateTable
CREATE TABLE "User" (
    "acct_num" VARCHAR(15) NOT NULL,
    "name" VARCHAR(20) NOT NULL,
    "email" VARCHAR(30) NOT NULL,
    "password" VARCHAR(200) NOT NULL,
    "address" VARCHAR(100) NOT NULL,

    CONSTRAINT "User_pkey" PRIMARY KEY ("acct_num")
);

-- CreateTable
CREATE TABLE "AccountDetail" (
    "loc_acct" VARCHAR(15) NOT NULL,
    "prin_pan" VARCHAR(50) NOT NULL,
    "acct_typ" VARCHAR(2) NOT NULL,
    "actv_typ" VARCHAR(1) NOT NULL,
    "blnc_amt" DOUBLE PRECISION NOT NULL,
    "loan_amt" DOUBLE PRECISION NOT NULL,
    "cycc_day" INTEGER NOT NULL,
    "min_loan_pymnt" DOUBLE PRECISION NOT NULL,
    "acct_num" VARCHAR(15) NOT NULL,

    CONSTRAINT "AccountDetail_pkey" PRIMARY KEY ("loc_acct")
);

-- CreateTable
CREATE TABLE "TransactionDetail" (
    "trx_id" VARCHAR(15) NOT NULL,
    "timestamps" TIMESTAMP(3) NOT NULL,
    "receiver_pan" VARCHAR(50) NOT NULL,
    "sender_pan" VARCHAR(50) NOT NULL,
    "apv_code" VARCHAR(10) NOT NULL,
    "trx_typ" VARCHAR(2) NOT NULL,
    "amt" DOUBLE PRECISION NOT NULL,
    "status" VARCHAR(1) NOT NULL,
    "desc" VARCHAR(100) NOT NULL,
    "loc_acct" VARCHAR(15) NOT NULL,

    CONSTRAINT "TransactionDetail_pkey" PRIMARY KEY ("trx_id")
);

-- CreateIndex
CREATE UNIQUE INDEX "User_acct_num_key" ON "User"("acct_num");

-- CreateIndex
CREATE UNIQUE INDEX "User_email_key" ON "User"("email");

-- CreateIndex
CREATE UNIQUE INDEX "AccountDetail_loc_acct_key" ON "AccountDetail"("loc_acct");

-- CreateIndex
CREATE UNIQUE INDEX "AccountDetail_prin_pan_key" ON "AccountDetail"("prin_pan");

-- CreateIndex
CREATE UNIQUE INDEX "TransactionDetail_trx_id_key" ON "TransactionDetail"("trx_id");

-- AddForeignKey
ALTER TABLE "AccountDetail" ADD CONSTRAINT "AccountDetail_acct_num_fkey" FOREIGN KEY ("acct_num") REFERENCES "User"("acct_num") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "TransactionDetail" ADD CONSTRAINT "TransactionDetail_loc_acct_fkey" FOREIGN KEY ("loc_acct") REFERENCES "AccountDetail"("loc_acct") ON DELETE RESTRICT ON UPDATE CASCADE;
