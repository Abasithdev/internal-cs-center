import api from "./axiosClient";
import type { PaymentResponse } from "@/type/payment";

export async function getPayments(params: Record<string, string | number | boolean> ): Promise<PaymentResponse> {
    const{data} = await api.get("/payments",{params});
    return data;
}

export async function reviewPayment(id:string): Promise<void> {
    await api.put(`/payments/${id}/review`);
}