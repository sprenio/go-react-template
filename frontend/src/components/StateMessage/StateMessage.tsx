import { useMessage } from '@/providers/MessageProvider';
import {MessageBox} from '@/components/MessageBox';

export default function StateMessage() {
  const { message } = useMessage();
  if (!message?.text) return null;

  return <MessageBox message={message.text} type={message.type}></MessageBox>;
}
