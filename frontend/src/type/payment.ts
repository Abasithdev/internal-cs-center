export interface Payment{
    id: string;
    merchant_name: string;
    date: string;
    amount:number;
    status: "completed" | "processing" | "failed";
    reviewed: boolean;
}

export interface PaymentMeta{
    total: number;
    total_pages: number;
    page: number;
    size: number;
    data: Payment[];
}

export interface PaymentSummary{
    total: number;
    completed: number;
    processing: number;
    failed: number;
}

export interface PaymentResponse{
    meta: PaymentMeta;
    summary: PaymentSummary;
}