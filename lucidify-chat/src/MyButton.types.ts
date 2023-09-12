import { ButtonHTMLAttributes } from "react";

export interface MyButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  size?: 'small' | 'medium' | 'large';
}
