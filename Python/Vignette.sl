class Vignette(color background = (0,0,0))
{
   public void imager(
          output varying color Ci;
          output varying color Oi;)
   {
        Ci = Ci + (1-Oi)*background;

        float x = s - 0.5;
        float y = t - 0.5;
        float d = x * x + y * y;
        Ci *= 1.0 - d;

        float n = (cellnoise(s*1000, t*1000) - 0.5);
        Ci += n * 0.025;
   }
}